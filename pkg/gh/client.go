package gh

import (
	"context"
	"errors"
	"github.com/google/go-github/v43/github"
	"github.com/masterminds/semver"
	"go.uber.org/zap"
	"time"
)

type VersionService struct {
	client   *github.Client
	versions map[string]*VersionInfo
	logger   *zap.Logger
	owner    string
	repo     string
}

type VersionInfo struct {
	SemVerVersion *semver.Version
	Name          string
}

func NewVersionService(logger *zap.Logger, owner, repo string) (*VersionService, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if len(owner) == 0 {
		return nil, errors.New("owner cannot be empty")
	}
	if len(repo) == 0 {
		return nil, errors.New("repo cannot be empty")
	}

	return &VersionService{
		client:   github.NewClient(nil),
		logger:   logger.Named("version_service"),
		owner:    owner,
		repo:     repo,
		versions: map[string]*VersionInfo{},
	}, nil
}

func (v *VersionService) GetVersions(ctx context.Context) (map[string]*VersionInfo, error) {
	if len(v.versions) == 0 {
		if err := v.RefreshVersions(ctx); err != nil {
			return nil, err
		}
	}
	return v.versions, nil
}

func (v *VersionService) GetVersion(ctx context.Context, ver *semver.Version) (*VersionInfo, bool) {
	if len(v.versions) == 0 {
		if err := v.RefreshVersions(ctx); err != nil {
			return nil, false
		}
	}
	retVer, exists := v.versions[ver.String()]
	return retVer, exists
}

func (v *VersionService) RefreshVersions(ctx context.Context) error {
	tags, err := v.getVersionTags(ctx)
	if err != nil {
		return err
	}
	for _, r := range tags {
		sv, err := semver.NewVersion(r)
		if err != nil {
			continue
		}
		v.versions[sv.String()] = &VersionInfo{
			SemVerVersion: sv,
			Name:          r,
		}
	}
	v.logger.Info("versions refresh complete")
	return nil
}

func (v *VersionService) StartRefresh(interval time.Duration, exit chan struct{}) {
	v.logger.Info("starting refresh process")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// tick
			if err := v.RefreshVersions(context.Background()); err != nil {
				v.logger.Error("unable to refresh tag cache", zap.Error(err))
			}
		case <-exit:
			v.logger.Info("refresh process exiting")
			return
		}
	}
}

func (v *VersionService) getVersionTags(ctx context.Context) ([]string, error) {
	var tags []string
	var releases []*github.RepositoryRelease
	listOpts := &github.ListOptions{
		PerPage: 100,
	}
	for {
		repoReleases, res, err := v.client.Repositories.ListReleases(ctx, v.owner, v.repo, listOpts)
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
	return tags, nil
}
