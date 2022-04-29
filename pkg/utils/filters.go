package utils

import (
	"fmt"
	"github.com/masterminds/semver"
)

func FilterVersions(versions []*semver.Version, major, minor string) ([]*semver.Version, error) {
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
