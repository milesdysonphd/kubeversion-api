package utils

import (
	"fmt"
	"github.com/masterminds/semver"
	"github.com/samber/lo"
	"github.com/superbased/kubeversion-api/pkg/gh"
	"golang.org/x/exp/maps"
	"sort"
)

func FilterVersions(versions map[string]*gh.VersionInfo, major, minor string) ([]*semver.Version, error) {
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
		// filter based on constraint
		retVersions = lo.FilterMap[*gh.VersionInfo, *semver.Version](maps.Values(versions), func(x *gh.VersionInfo, _ int) (*semver.Version, bool) {
			if semverConstraint.Check(x.SemVerVersion) {
				return x.SemVerVersion, true
			}
			return nil, false
		})
	} else {
		retVersions = lo.Map[*gh.VersionInfo, *semver.Version](maps.Values(versions), func(x *gh.VersionInfo, _ int) *semver.Version {
			return x.SemVerVersion
		})
	}
	sort.Sort(semver.Collection(retVersions))
	return retVersions, nil
}
