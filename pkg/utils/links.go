package utils

import (
	"fmt"
	"github.com/masterminds/semver"
)

const (
	baseRepo      = "https://github.com/kubernetes/kubernetes"
	baseShortUrl  = "https://git.k8s.io/kubernetes"
	changelogPath = "CHANGELOG/CHANGELOG"
)

var (
	bins = []string{
		"kubectl",
	}
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
	Version   string                    `json:"version"`
	Changelog string                    `json:"changelog"`
	Downloads map[string][]DownloadLink `json:"downloads"`
}

type DownloadLink struct {
	Binary       string `json:"binary"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	URL          string `json:"url"`
	ChecksumURL  string `json:"checksumUrl"`
}

func BuildVersionsResponse(versions []*semver.Version) []Release {
	releases := make([]Release, len(versions))
	for _, v := range versions {
		tmp := Release{
			Version:   v.String(),
			Changelog: fmt.Sprintf("%s/%s-%v.%v.md", baseShortUrl, changelogPath, v.Major(), v.Minor()),
			Downloads: map[string][]DownloadLink{},
		}
		for k, p := range platformArch {
			for _, b := range bins {
				binName := b
				if k == "windows" {
					binName = fmt.Sprintf("%s.exe", b)
				}

				for _, arch := range p {
					tmp.Downloads[k] = append(tmp.Downloads[k], DownloadLink{
						Binary:       "kubectl",
						Platform:     k,
						Architecture: arch,
						URL:          fmt.Sprintf("https://dl.k8s.io/v%s/bin/%s/%s/%s", v.String(), k, arch, binName),
						ChecksumURL:  fmt.Sprintf("https://dl.k8s.io/v%s/bin/%s/%s/%s.sha256", v.String(), k, arch, binName),
					})
				}
			}
		}
		releases = append(releases, tmp)
	}
	return releases
}

func BuildVersionResponse(v *semver.Version) Release {
	ret := Release{
		Version:   v.String(),
		Changelog: fmt.Sprintf("%s/%s-%v.%v.md", baseShortUrl, changelogPath, v.Major(), v.Minor()),
		Downloads: map[string][]DownloadLink{},
	}
	for k, p := range platformArch {
		for _, b := range bins {
			binName := b
			if k == "windows" {
				binName = fmt.Sprintf("%s.exe", b)
			}

			for _, arch := range p {
				ret.Downloads[k] = append(ret.Downloads[k], DownloadLink{
					Binary:       "kubectl",
					Platform:     k,
					Architecture: arch,
					URL:          fmt.Sprintf("https://dl.k8s.io/v%s/bin/%s/%s/%s", v.String(), k, arch, binName),
					ChecksumURL:  fmt.Sprintf("https://dl.k8s.io/v%s/bin/%s/%s/%s.sha256", v.String(), k, arch, binName),
				})
			}
		}
	}
	return ret
}
