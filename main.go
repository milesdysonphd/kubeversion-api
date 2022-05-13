package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	v1 "github.com/superbased/kubeversion-api/controllers/v1"
	"github.com/superbased/kubeversion-api/pkg/gh"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// logging
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := conf.Build()
	if err != nil {
		log.Fatal(err)
	}

	// echo framework
	e := echo.New()

	// new relic
	if len(os.Getenv("LOCAL_MODE")) == 0 {
		newRelicKey, ok := os.LookupEnv("NEW_RELIC_KEY")
		if !ok {
			logger.Fatal("NEW_RELIC_KEY environment variable is missing")
		}
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("KubeVersion API"),
			newrelic.ConfigLicense(newRelicKey),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			logger.Fatal("error instantiating new relic", zap.Error(err))
		}
		e.Use(nrecho.Middleware(app))
	}

	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// exit channel
	exit := make(chan struct{}, 1)
	versionService, err := gh.NewVersionService(logger, "kubernetes", "kubernetes")

	// start the version refresh service
	logger.Info("starting versions refresh")
	if err := versionService.RefreshVersions(context.Background()); err != nil {
		logger.Fatal("unable to fetch initial version data", zap.Error(err))
	}
	go versionService.StartRefresh(1*time.Hour, exit)

	v1VersionsController, err := v1.NewVersionsController(versionService, logger)
	if err != nil {
		logger.Fatal("unable to build versions controller", zap.Error(err))
	}

	// mount /v1/versions endpoints
	v1VersionsController.Mount(e)

	// mount a default endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "hello there",
		})
	})

	// mount a health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	if err := e.Start(":8080"); err != nil {
		logger.Error("error running server", zap.Error(err))
	}
	close(exit)
}
