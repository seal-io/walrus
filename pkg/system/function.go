package system

import (
	"github.com/seal-io/utils/varx"
	"k8s.io/apimachinery/pkg/util/sets"
)

// DisableApplications is a set of applications that are not allowed to be installed.
var DisableApplications = varx.NewOnce(sets.New[string]())

// ConfigureDisallowApplications configures the applications of the system which not be installed.
func ConfigureDisallowApplications(items []string) {
	DisableApplications.Configure(sets.New[string](items...))
}
