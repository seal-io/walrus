package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/seal-io/walrus/utils/strs"
)

var (
	Version      = "dev"
	GitCommit    = "HEAD"
	instanceUUID = strs.Hex(16)
)

func Get() string {
	return fmt.Sprintf("%s (%s)", Version, GitCommit)
}

func GetUserAgent() string {
	return GetUserAgentWith("walrus")
}

func GetUserAgentWith(name string) string {
	return "seal.io/" + name + "; version=" + Get() + "; os=" + runtime.GOOS + "; arch=" + runtime.GOARCH
}

func GetInstanceUUID() string {
	return instanceUUID
}

func Major() string {
	vX := semver.Major(Version)
	if vX == "" {
		return Version
	}

	return vX
}

func MajorMinor() string {
	vXy := semver.MajorMinor(Version)
	if vXy == "" {
		return Version
	}

	return vXy
}

func Previous() string {
	vXy := MajorMinor()
	if vXy == Version {
		return Version
	}

	v := strings.Split(vXy, ".")
	if v[1] != "0" {
		y, _ := strconv.ParseInt(v[1], 10, 64)
		y--

		if y >= 0 {
			return v[0] + "." + strconv.FormatInt(y, 10)
		}
	}

	x, _ := strconv.ParseInt(v[0][1:], 10, 64)
	x--

	if x < 0 {
		return Version
	}

	return "v" + strconv.FormatInt(x, 10)
}

func IsValid() bool {
	return semver.IsValid(Version)
}

func IsDevVersion() bool {
	return Version == "" ||
		Version == "dev" ||
		Version == "main"
}
