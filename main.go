package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	v1 "github.com/superbased/kubeversion-api/controllers/v1"
	"github.com/superbased/kubeversion-api/pkg/gh"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
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

	// fiber framework
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// middleware
	app.Use(requestid.New())

	// exit channel
	exit := make(chan struct{}, 1)
	versionService, err := gh.NewVersionService(logger, "kubernetes", "kubernetes")
	// start the version refresh service
	logger.Info("starting version refresh service", zap.Duration("interval", 1*time.Hour))
	if err := versionService.RefreshVersions(context.Background()); err != nil {
		logger.Fatal("unable to fetch initial version data", zap.Error(err))
	}
	go versionService.StartRefresh(1*time.Hour, exit)

	v1VersionsController, err := v1.NewVersionsController(versionService, logger)
	if err != nil {
		logger.Fatal("unable to build versions controller", zap.Error(err))
	}

	// mount /v1/versions endpoints
	v1VersionsController.Mount(app)

	// mount a default endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "hello there",
		})
	})

	if err := app.Listen(":8080"); err != nil {
		logger.Error("error running server", zap.Error(err))
	}
	close(exit)
}
