package gh

import (
	"context"
	"errors"
	"github.com/google/go-github/v43/github"
	"github.com/masterminds/semver"
	"go.uber.org/zap"
	"sort"
	"time"
)

type VersionService struct {
	client   *github.Client
	versions []*semver.Version
	logger   *zap.Logger
	owner    string
	repo     string
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
		client: github.NewClient(nil),
		logger: logger.Named("version_service"),
		owner:  owner,
		repo:   repo,
	}, nil
}

func (v *VersionService) GetVersions(ctx context.Context) ([]*semver.Version, error) {
	if len(v.versions) == 0 {
		if err := v.RefreshVersions(ctx); err != nil {
			return nil, err
		}
	}
	return v.versions, nil
}

func (v *VersionService) RefreshVersions(ctx context.Context) error {
	tags, err := v.getVersionTags(ctx)
	if err != nil {
		return err
	}
	vs := make([]*semver.Version, len(tags))
	for i, r := range tags {
		v, err := semver.NewVersion(r)
		if err != nil {
			continue
		}
		vs[i] = v
	}
	sort.Sort(semver.Collection(vs))
	v.versions = vs
	return nil
}

func (v *VersionService) StartRefresh(interval time.Duration, exit chan struct{}) {
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
