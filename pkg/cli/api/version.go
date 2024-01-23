package api

import (
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/version"
)

// VersionInfo include the client, server and openapi version.
type VersionInfo struct {
	ClientVersion  Version `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`
	ServerVersion  Version `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
	OpenAPIVersion Version `json:"openAPIVersion,omitempty" yaml:"-"`
}

// Version include the version and commit.
type Version struct {
	Version      string `json:"version" yaml:"version"`
	GitCommit    string `json:"gitCommit" yaml:"gitCommit"`
	IsDevVersion bool   `json:"isDevVersion" yaml:"isDevVersion"`
}

// GetVersion get client, server and openapi version.
func GetVersion(sc *config.Config) (*VersionInfo, error) {
	// Client version.
	info := &VersionInfo{
		ClientVersion: Version{
			Version:      version.Version,
			GitCommit:    version.GitCommit,
			IsDevVersion: version.IsDevVersion(),
		},
	}

	if sc != nil && sc.Reachable {
		// Server version.
		sv, err := sc.ServerVersion()
		if err != nil {
			return nil, err
		}

		info.ServerVersion = Version{
			Version:      sv.Version,
			GitCommit:    sv.Commit,
			IsDevVersion: version.IsDevVersionWith(sv.Version),
		}

		// Openapi version.
		av := GetAPIVersionFromCache()
		if av != nil {
			info.OpenAPIVersion = Version{
				Version:      av.Version,
				GitCommit:    av.GitCommit,
				IsDevVersion: version.IsDevVersionWith(av.Version),
			}
		}
	}

	return info, nil
}
