package openapi

components: {
	// Schemas are injected from the types.cue file
	schemas: {}
}

paths: {
	"/": get: {
		summary:     "Welcome Message"
		description: "A convinient endpoint for your sanity test."
		operationId: "WelcomeGet"
		responses: "200": {
			description: "Successful Response"
			content: "application/json": schema: $ref: "#/components/schemas/Welcome"
		}
	}
	"/config": get: {
		summary:     "View running config"
		description: "See the current config"
		operationId: "ConfigGet"
		responses: "200": {
			description: "Successful Response"
			content: "application/json": schema: $ref: "#/components/schemas/Config"
		}
	}
	"/metrics/pricing/azure": get: {
		summary:     "Azure Pricing"
		description: "A Prometheus metrics formatted overview of Azure Prices"
		operationId: "MetricsPricingAzureGet"
		responses: "200": {
			description: "Successful Response"
			content: "text/plain": schema: {}
		}
	}
}
