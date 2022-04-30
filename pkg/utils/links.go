package utils

import (
	"fmt"
	"github.com/masterminds/semver"
)

var (
	platformArch = map[string][]string{
		"darwin": {
			"amd64",
			"arm64",
		},
		"windows": {
			"386",
			"amd64",
		},
		"linux": {
			"386",
			"amd64",
			"arm",
			"arm64",
			"ppc64le",
			"s390x",
		},
	}
)

type Release struct {
	Version   string         `json:"version"`
	Downloads []DownloadLink `json:"downloads"`
}

type DownloadLink struct {
	Binary       string `json:"binary"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	URL          string `json:"url"`
	ChecksumURL  string `json:"checksumUrl"`
}

func BuildVersionResponse(versions []*semver.Version) []Release {
	releases := make([]Release, len(versions))
	for _, v := range versions {
		tmp := Release{
			Version:   v.String(),
			Downloads: []DownloadLink{},
		}
		for k, p := range platformArch {
			binName := "kubectl"
			if k == "windows" {
				binName = "kubectl.exe"
			}
			for _, arch := range p {
				tmp.Downloads = append(tmp.Downloads, DownloadLink{
					Binary:       "kubectl",
					Platform:     k,
					Architecture: arch,
					URL:          fmt.Sprintf("https://dl.k8s.io/v%s/bin/%s/%s/%s", v.String(), k, arch, binName),
					ChecksumURL:  fmt.Sprintf("https://dl.k8s.io/v%s/bin/%s/%s/%s.sha256", v.String(), k, arch, binName),
				})
			}
		}
		releases = append(releases, tmp)
	}
	return releases
}
