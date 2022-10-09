package openapi

components: {
	// Schemas are injected from the types.cue file
	schemas: {}
}

paths: "/": get: {
	summary:     "Welcome Message"
	description: "A convinient endpoint for your sanity test."
	operationId: "WelcomeGet"
	responses: "200": {
		description: "Successful Response"
		content: "text/plain": schema: {}
	}
}
