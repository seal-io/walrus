package api

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/utils/log"
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

	if sc != nil && sc.Server != "" && sc.Token != "" {
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

// CompareVersion compare the client, server and openapi version.
func CompareVersion(sc *config.Config) (shouldUpdate bool, err error) {
	err = sc.ValidateAndSetup()
	if err != nil {
		return false, err
	}

	v, err := GetVersion(sc)
	if err != nil {
		return false, err
	}

	switch {
	case !v.ClientVersion.IsDevVersion && !v.ServerVersion.IsDevVersion:
		// Release version, check if client and server version match.
		cv, err := semver.NewVersion(v.ClientVersion.Version)
		if err != nil {
			return false, err
		}

		sv, err := semver.NewVersion(v.ServerVersion.Version)
		if err != nil {
			return false, err
		}

		if cv.Major() != sv.Major() {
			return false,
				fmt.Errorf("major version incompatibility detected between cli (%s) and server (%s)", cv, sv)
		}

		if cv.Minor() != sv.Minor() {
			log.Warnf(
				// nolint:lll
				"minor version mismatch detected between cli (%s) and server (%s). It is recommended to use a compatible version",
				v.ClientVersion.Version,
				v.ServerVersion.Version,
			)
		}

	case v.ClientVersion.IsDevVersion && !v.ServerVersion.IsDevVersion:
		// Client is release version, server is dev version.
		log.Warnf("cli (%s) is running the development version, whereas the server (%s) is using a release version",
			v.ClientVersion.Version,
			v.ServerVersion.Version)
	}

	return v.ServerVersion.GitCommit != v.OpenAPIVersion.GitCommit, nil
}
