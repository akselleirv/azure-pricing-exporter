package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/akselleirv/azure-pricing-exporter/api/models"
	"github.com/akselleirv/azure-pricing-exporter/internal/azure"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	currencyCode    = "currencyCode"
	location        = "location"
	armSkuName      = "armSkuName"
	reservationTerm = "reservationTerm"
)

type Server struct {
	reg      *prometheus.Registry
	resolver azure.AzurePriceResolver
	cfg      models.Config
}

var priceGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "azure_prices",
		Help: "Azure prices",
	},
	[]string{currencyCode, location, armSkuName, reservationTerm},
)

func NewCostsService(configPath string, resolver azure.AzurePriceResolver, done <-chan struct{}) (*Server, error) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(priceGauge)

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg := models.Config{}
	if err = json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	validator := validator.New()
	cfg.Timestamp = time.Now().String()

	if err = validator.Struct(cfg); err != nil {
		return nil, err
	}

	s := Server{reg: reg, resolver: resolver, cfg: cfg}
	s.collectPrices(time.Duration(cfg.IntervalInMinutes)*time.Minute, cfg.ConcurrencyLevel, done)
	return &s, nil
}

func (s *Server) collectPrices(interval time.Duration, numOfWorkers int, done <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		jobsCh := make(chan models.AzurePriceResolverOpts, numOfWorkers)
		resultCh := make(chan azure.AzurePriceItems, numOfWorkers)

		for i := 0; i < numOfWorkers; i++ {
			go s.resolver.Resolve(jobsCh, resultCh, done)
		}

		for _, resource := range s.cfg.ResolvePricesFor {
			jobsCh <- resource
		}

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, resource := range s.cfg.ResolvePricesFor {
					jobsCh <- resource
				}
			case resp := <-resultCh:
				if resp.Err != nil {
					log.Println(resp.Err)
					continue
				}
				for _, p := range resp.Items {
					retailPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", p.RetailPrice), 64)
					priceGauge.With(
						prometheus.Labels{
							currencyCode:    p.CurrencyCode,
							location:        p.Location,
							armSkuName:      p.ArmSkuName,
							reservationTerm: p.ReservationTerm,
						},
					).Set(retailPrice)
				}
			}
		}
	}()
}

func (s *Server) WelcomeGet(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Welcome{Message: "hello from azure-pricing-exporter"})
}

func (s *Server) MetricsPricingAzureGet(c echo.Context) error {
	promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{}).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) ConfigGet(c echo.Context) error {
	return c.JSON(http.StatusOK, s.cfg)
}
