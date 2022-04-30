package v1

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/superbased/kubeversion-api/pkg/gh"
	"github.com/superbased/kubeversion-api/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// VersionsController handles all /versions endpoints
type VersionsController struct {
	logger         *zap.Logger
	versionService *gh.VersionService
}

// NewVersionsController returns a new VersionsController
func NewVersionsController(versionService *gh.VersionService, logger *zap.Logger) (*VersionsController, error) {
	return &VersionsController{
		logger:         logger.Named("versions_controller").WithOptions(zap.Fields(zap.String("controller_version", "v1"))),
		versionService: versionService,
	}, nil
}

// Mount is responsible for mounting the /versions endpoints
func (v *VersionsController) Mount(a *fiber.App) {
	v1Endpoints := a.Group("/v1/versions")
	v1Endpoints.Get("/", v.getVersions)
	v1Endpoints.Get("/latest", v.getLatestVersion)
}

func (v *VersionsController) getVersions(c *fiber.Ctx) error {
	major := c.Query("major")
	minor := c.Query("minor")

	ctx, cancel := context.WithTimeout(c.Context(), time.Second*10)
	defer cancel()
	versions, err := v.versionService.GetVersions(ctx)
	if err != nil {
		v.logger.Error("error getting versions", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]string{
			"message": "something went wrong",
		})
	}

	retVersions, err := utils.FilterVersions(versions, major, minor)
	if err != nil {
		v.logger.Error("error filtering versions", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return c.JSON(map[string]string{
			"message": "unable to build version filter",
		})
	}

	return c.JSON(map[string]interface{}{
		"data": map[string]interface{}{
			"versions": retVersions,
		},
	})
}

func (v *VersionsController) getLatestVersion(c *fiber.Ctx) error {
	major := c.Query("major")
	minor := c.Query("minor")

	ctx, cancel := context.WithTimeout(c.Context(), time.Second*10)
	defer cancel()
	versions, err := v.versionService.GetVersions(ctx)
	if err != nil {
		v.logger.Error("error getting versions", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return c.JSON(map[string]string{
			"message": "something went wrong",
		})
	}

	retVersions, err := utils.FilterVersions(versions, major, minor)
	if err != nil {
		v.logger.Error("error filtering versions", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return c.JSON(map[string]string{
			"message": "unable to build version filter",
		})
	}

	releases := utils.BuildVersionResponse(retVersions)
	return c.JSON(map[string]interface{}{
		"data": releases[len(releases)-1],
	})
}
