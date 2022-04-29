package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/masterminds/semver"
	"github.com/milesdysonphd/kubeversion-api/pkg/gh"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
)

func main() {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := conf.Build()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(requestid.New())

	exit := make(chan struct{}, 1)
	versionServer, err := gh.NewVersionService(logger, "kubernetes", "kubernetes")

	logger.Info("starting version refresh service", zap.Duration("interval", 1*time.Hour))
	if err := versionServer.RefreshVersions(context.Background()); err != nil {
		logger.Fatal("unable to fetch initial version data", zap.Error(err))
	}
	go versionServer.StartRefresh(1*time.Hour, exit)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "hello there",
		})
	})

	app.Get("/versions", func(c *fiber.Ctx) error {
		major := c.Query("major")
		minor := c.Query("minor")

		ctx, cancel := context.WithTimeout(c.Context(), time.Second*10)
		defer cancel()
		versions, err := versionServer.GetVersions(ctx)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(map[string]string{
				"message": "something went wrong",
			})
		}

		retVersions, err := filterVersions(versions, major, minor)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(map[string]string{
				"message": "unable to build version filter",
			})
		}

		return c.JSON(map[string]interface{}{
			"versions": retVersions,
		})
	})

	app.Get("/versions/latest", func(c *fiber.Ctx) error {
		major := c.Query("major")
		minor := c.Query("minor")

		ctx, cancel := context.WithTimeout(c.Context(), time.Second*10)
		defer cancel()
		versions, err := versionServer.GetVersions(ctx)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(map[string]string{
				"message": "something went wrong",
			})
		}

		retVersions, err := filterVersions(versions, major, minor)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(map[string]string{
				"message": "unable to build version filter",
			})
		}

		return c.JSON(map[string]interface{}{
			"version": retVersions[len(retVersions)-1],
		})
	})

	if err := app.Listen(":8080"); err != nil {
		logger.Error("error running server", zap.Error(err))
	}
	close(exit)
}

func filterVersions(versions []*semver.Version, major, minor string) ([]*semver.Version, error) {
	var retVersions []*semver.Version
	if len(major) != 0 {
		filterConstraint := fmt.Sprintf("~%s", major)
		if len(minor) != 0 {
			filterConstraint = fmt.Sprintf("~%s.%s", major, minor)
		}
		semverConstraint, err := semver.NewConstraint(filterConstraint)
		if err != nil {
			return nil, err
		}
		for _, v := range versions {
			if semverConstraint.Check(v) {
				retVersions = append(retVersions, v)
			}
		}
	} else {
		retVersions = versions
	}
	return retVersions, nil
}
