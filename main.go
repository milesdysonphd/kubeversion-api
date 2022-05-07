package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// sentry
	sentryDSN, ok := os.LookupEnv("SENTRY_DSN")
	if !ok {
		logger.Fatal("SENTRY_DSN environment variable is missing")
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		TracesSampleRate: 0.2,
	}); err != nil {
		logger.Fatal("unable to initialize sentry", zap.Error(err))
	}
	defer sentry.Flush(2 * time.Second)

	// echo framework
	e := echo.New()
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
