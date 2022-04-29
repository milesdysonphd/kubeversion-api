package utils

import (
	"github.com/masterminds/semver"
	"reflect"
	"testing"
)

func TestFilterVersions(t *testing.T) {
	type args struct {
		versions []*semver.Version
		major    string
		minor    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*semver.Version
		wantErr bool
	}{
		{
			name: "filter major only",
			args: args{
				versions: generateVersions("v1.3.5", "v1.0.0", "v0.0.1", "v1.0.5"),
				major:    "1",
				minor:    "",
			},
			want:    generateVersions("v1.3.5", "v1.0.0", "v1.0.5"),
			wantErr: false,
		},
		{
			name: "filter major and minor",
			args: args{
				versions: generateVersions("v1.3.5", "v1.0.0", "v0.0.1", "v1.0.5"),
				major:    "1",
				minor:    "0",
			},
			want:    generateVersions("v1.0.0", "v1.0.5"),
			wantErr: false,
		},
		{
			name: "filter minor only default to no filters",
			args: args{
				versions: generateVersions("v1.3.5", "v1.0.0", "v0.0.1", "v1.0.5"),
				major:    "",
				minor:    "1",
			},
			want:    generateVersions("v1.3.5", "v1.0.0", "v0.0.1", "v1.0.5"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterVersions(tt.args.versions, tt.args.major, tt.args.minor)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterVersions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateVersions(versions ...string) []*semver.Version {
	var pVersions []*semver.Version
	for _, v := range versions {
		sVersion, _ := semver.NewVersion(v)
		pVersions = append(pVersions, sVersion)
	}
	return pVersions
}
