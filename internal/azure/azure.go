package azure

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/akselleirv/azure-pricing-exporter/api/models"
	"github.com/go-playground/validator/v10"
)

type AzurePrice struct {
	c        *http.Client
	validate *validator.Validate
}

type AzurePriceResolver interface {
	Resolve(jobsCh <-chan models.AzurePriceResolverOpts, resultCh chan<- AzurePriceItems, done <-chan struct{})
}

func New() *AzurePrice {
	return &AzurePrice{c: &http.Client{Timeout: 5 * time.Second}, validate: validator.New()}
}

// See docs: https://learn.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices#api-property-details
type AzurePriceItem struct {
	// Azure data center where the resource is deployed
	Location string `json:"location" validate:"required"`

	// The currency in which rates are defined.
	CurrencyCode string `json:"currencyCode" validate:"required"`

	// SKU name registered in Azure. E.g. Standard_B2ms
	ArmSkuName string `json:"armSkuName" validate:"required"`

	ProductName string `json:"productName" validate:"required"`

	// Prices per hour without discount
	RetailPrice float32 `json:"retailPrice" validate:"required"`

	// Pay-as-you-go, one year or three years
	ReservationTerm string `json:"reservationTerm"`
}

type AzurePriceItems struct {
	Items []AzurePriceItem
	Err   error
}

func (a *AzurePrice) Resolve(jobsCh <-chan models.AzurePriceResolverOpts, resultCh chan<- AzurePriceItems, done <-chan struct{}) {
	for {
		select {
		case job := <-jobsCh:
			log.Printf("[INFO] collecting prices for '%s' in '%s' with currency '%s'", job.ArmSkuName, job.Location, job.CurrencyCode)
			items, err := a.doResolve(&job)
			if err != nil {
				resultCh <- AzurePriceItems{
					Err: err,
				}
				continue
			}
			if len(items) == 0 {
				resultCh <- AzurePriceItems{
					Err: fmt.Errorf("[WARNING] unable to find prices for '%s' for '%s' in currency '%s'", job.ArmSkuName, job.Location, job.CurrencyCode),
				}
				continue
			}
			resultCh <- AzurePriceItems{
				Err:   err,
				Items: items,
			}
		case <-done:
			return
		}
	}

}

func (a *AzurePrice) doResolve(opts *models.AzurePriceResolverOpts) ([]AzurePriceItem, error) {
	if err := a.validate.Struct(opts); err != nil {
		return nil, err

	}
	req, err := http.NewRequest(http.MethodGet, "https://prices.azure.com/api/retail/prices", nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Add("api-version", "2021-10-01-preview")
	queryParams.Add("currencyCode", opts.CurrencyCode)
	queryParams.Add("$filter", buildFilterQuery(opts))

	req.URL.RawQuery = queryParams.Encode()
	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: got status code '%d'", string(body), resp.StatusCode)
	}
	type response struct {
		Items []AzurePriceItem `json:"Items"`
	}
	unmarshaledBody := response{}
	if err = json.Unmarshal(body, &unmarshaledBody); err != nil {
		return nil, err
	}
	var result []AzurePriceItem
	for _, item := range unmarshaledBody.Items {
		if err = a.validate.Struct(item); err != nil {
			return nil, fmt.Errorf("validation of response failed: %w", err)
		}
		if item.ReservationTerm == "" {
			item.ReservationTerm = "pay-as-you-go"
		}
		if strings.Contains(item.ProductName, "windows") {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}

func buildFilterQuery(opts *models.AzurePriceResolverOpts) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("armRegionName eq '%s' ", opts.Location))
	b.WriteString(fmt.Sprintf("and armSkuName eq '%s' ", opts.ArmSkuName))
	return b.String()
}
