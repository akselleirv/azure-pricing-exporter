package main

import (
	"log"
	"os"

	"github.com/akselleirv/azure-pricing-exporter/api"
	"github.com/akselleirv/azure-pricing-exporter/internal/azure"
	oapiCodegenMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/heptiolabs/healthcheck"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalln("error loading swagger spec: ", err.Error())
	}
	resolver := azure.New()
	done := make(chan struct{})

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.json"
	}

	service, err := api.NewCostsService(cfgPath, resolver, done)
	if err != nil {
		log.Fatalln("error creating service", err.Error())
	}
	e := echo.New()
	healthProbes(e)
	setupMetrics(e)
	e.Use(echoMiddleware.Recover())

	// If the apiRouter is not created, then the OpenAPI request validator will
	// block any requests which are not mentioned in the spec.
	apiRouter := e.Group("")
	apiRouter.Use(oapiCodegenMiddleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(apiRouter, service)
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":8080"))
}

// setupMetrics starts a metrics server on :9090
func setupMetrics(e *echo.Echo) {
	prom := prometheus.NewPrometheus("", nil)
	e.Use(prom.HandlerFunc)
	metricServer := echo.New()
	metricServer.HidePort = true
	metricServer.HideBanner = true
	prom.SetMetricsPath(metricServer)

	go func() {
		e.Logger.Fatal(metricServer.Start(":9090"))
	}()
}

func healthProbes(e *echo.Echo) {
	health := healthcheck.NewHandler()
	health.AddReadinessCheck("todo", func() error {
		return nil
	})

	e.GET("/ready", func(c echo.Context) error {
		health.ReadyEndpoint(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/live", func(c echo.Context) error {
		health.LiveEndpoint(c.Response().Writer, c.Request())
		return nil
	})
}
