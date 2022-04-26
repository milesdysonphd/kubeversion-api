package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v43/github"
	"github.com/labstack/echo/v4"
	"github.com/masterminds/semver"
	"log"
	"net/http"
	"os"
	"sort"
)

var (
	cacheVersionStrings []string
	cacheParsedVersions []*semver.Version
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	e := echo.New()
	e.GET("/", helloWorld)
	e.GET("/versions", listVersions)
	e.GET("/versions/latest", getLatestVersion)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func helloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "hello world",
	})
}

func listVersions(c echo.Context) error {
	if len(cacheVersionStrings) == 0 {
		versions, err := getReleases(context.TODO(), "kubernetes", "kubernetes")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "something went wrong",
			})
		}
		var retVersions []string
		for _, v := range versions {
			retVersions = append(retVersions, v.String())
		}
		cacheVersionStrings = retVersions
		cacheParsedVersions = versions
	}
	major := c.QueryParam("major")
	minor := c.QueryParam("minor")

	if len(major) > 0 {
		versionConstraintString := fmt.Sprintf("~%v", major)
		if len(minor) > 0 {
			versionConstraintString = fmt.Sprintf("~%v.%v", major, minor)
		}
		versionConstraint, err := semver.NewConstraint(versionConstraintString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "there was an error parsing your version constraints",
			})
		}
		var filteredVersions []string
		for _, v := range cacheParsedVersions {
			if versionConstraint.Check(v) {
				filteredVersions = append(filteredVersions, v.String())
			}
		}
		return c.JSON(http.StatusOK, map[string][]string{
			"versions": filteredVersions,
		})
	}

	return c.JSON(http.StatusOK, map[string][]string{
		"versions": cacheVersionStrings,
	})
}

func getLatestVersion(c echo.Context) error {
	if len(cacheVersionStrings) == 0 {
		versions, err := getReleases(context.TODO(), "kubernetes", "kubernetes")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "something went wrong",
			})
		}
		var retVersions []string
		for _, v := range versions {
			retVersions = append(retVersions, v.String())
		}
		cacheVersionStrings = retVersions
		cacheParsedVersions = versions
	}

	major := c.QueryParam("major")
	minor := c.QueryParam("minor")

	if len(major) > 0 {
		versionConstraintString := fmt.Sprintf("~%v", major)
		if len(minor) > 0 {
			versionConstraintString = fmt.Sprintf("~%v.%v", major, minor)
		}
		versionConstraint, err := semver.NewConstraint(versionConstraintString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "there was an error parsing your version constraints",
			})
		}
		var filteredVersions []string
		for _, v := range cacheParsedVersions {
			if versionConstraint.Check(v) {
				filteredVersions = append(filteredVersions, v.String())
			}
		}
		return c.JSON(http.StatusOK, map[string]string{
			"version": filteredVersions[len(filteredVersions)-1],
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"version": cacheVersionStrings[len(cacheVersionStrings)-1],
	})
}

func getReleases(ctx context.Context, owner, repo string) ([]*semver.Version, error) {
	client := github.NewClient(nil)
	var tags []string
	var releases []*github.RepositoryRelease
	listOpts := &github.ListOptions{
		PerPage: 100,
	}
	for {
		repoReleases, res, err := client.Repositories.ListReleases(ctx, owner, repo, listOpts)
		if err != nil {
			return nil, err
		}
		releases = append(releases, repoReleases...)
		if res.NextPage == 0 {
			break
		}
		listOpts.Page = res.NextPage
	}
	for _, v := range releases {
		if !*v.Prerelease {
			tags = append(tags, *v.TagName)
		}
	}

	vs := make([]*semver.Version, len(tags))
	for i, r := range tags {
		v, err := semver.NewVersion(r)
		if err != nil {
			fmt.Printf("error parsing %s\n", r)
			continue
		}
		vs[i] = v
	}
	sort.Sort(semver.Collection(vs))

	return vs, nil
}
