// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package models

// AzurePriceResolverOpts defines model for AzurePriceResolverOpts.
type AzurePriceResolverOpts struct {
	ArmSkuName   string `json:"armSkuName"`
	CurrencyCode string `json:"currencyCode"`
	Location     string `json:"location"`
}

// Config defines model for Config.
type Config struct {
	ConcurrencyLevel  int                      `json:"concurrencyLevel"`
	IntervalInMinutes float32                  `json:"intervalInMinutes"`
	ResolvePricesFor  []AzurePriceResolverOpts `json:"resolvePricesFor"`
	Timestamp         string                   `json:"timestamp"`
}

// Welcome defines model for Welcome.
type Welcome struct {
	// Message is a welcome message
	Message string `json:"message"`
}