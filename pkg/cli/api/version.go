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
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"gitCommit" yaml:"gitCommit"`
}

// GetVersion get client, server and openapi version.
func GetVersion(sc *config.Config) *VersionInfo {
	// Client version.
	info := &VersionInfo{
		ClientVersion: Version{
			Version:   version.Version,
			GitCommit: version.GitCommit,
		},
	}

	if sc != nil && sc.Reachable {
		// Server version.
		sv, err := sc.ServerVersion()
		if err != nil {
			// Return client version if server version is not reachable.
			return info
		}

		info.ServerVersion = Version{
			Version:   sv.Version,
			GitCommit: sv.Commit,
		}

		// Openapi version.
		av := GetAPIVersionFromCache()
		if av != nil {
			info.OpenAPIVersion = Version{
				Version:   av.Version,
				GitCommit: av.GitCommit,
			}
		}
	}

	return info
}
