package v1

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/masterminds/semver"
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
func (v *VersionsController) Mount(e *echo.Echo) {
	g := e.Group("/v1/versions")
	g.GET("", v.getVersions)
	g.GET("/:version", v.getVersion)
	g.GET("/latest", v.getLatestVersion)
}

func (v *VersionsController) getVersions(c echo.Context) error {
	major := c.QueryParam("major")
	minor := c.QueryParam("minor")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	versions, err := v.versionService.GetVersions(ctx)
	if err != nil {
		v.logger.Error("error getting versions", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "something went wrong",
		})
	}

	retVersions, err := utils.FilterVersions(versions, major, minor)
	if err != nil {
		v.logger.Error("error filtering versions", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "unable to build version filter",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"versions": retVersions,
		},
	})
}

func (v *VersionsController) getVersion(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	version := c.Param("version")
	parsedVersion, err := semver.NewVersion(version)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "unable to parse provided version",
		})
	}
	retVersion, exists := v.versionService.GetVersion(ctx, parsedVersion)
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "unable to find specified version",
		})
	}

	ret := utils.BuildVersionResponse(retVersion.SemVerVersion)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": ret,
	})
}

func (v *VersionsController) getLatestVersion(c echo.Context) error {
	major := c.QueryParam("major")
	minor := c.QueryParam("minor")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	versions, err := v.versionService.GetVersions(ctx)
	if err != nil {
		v.logger.Error("error getting versions", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "something went wrong",
		})
	}

	retVersions, err := utils.FilterVersions(versions, major, minor)
	if err != nil {
		v.logger.Error("error filtering versions", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "unable to build version filter",
		})
	}

	releases := utils.BuildVersionsResponse(retVersions)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": releases[len(releases)-1],
	})
}
