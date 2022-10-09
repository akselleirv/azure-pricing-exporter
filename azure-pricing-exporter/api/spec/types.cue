package openapi

info: {
	title:       "azure-pricing-exporter"
	description: ""
	version:     "0.0.1"
}

#Welcome: {
	// Message is a welcome message
	message: string
}

#Config: {
	intervalInMinutes: number
	concurrencyLevel:  int
	timestamp:         string
	resolvePricesFor: [...#AzurePriceResolverOpts]
}

#AzurePriceResolverOpts: {
	armSkuName:   string
	location:     string
	currencyCode: string
}
