package servervars

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/vars"
)

var (
	// NonLoopBackIPs is a set of non-loopback address of the host where the server is located.
	NonLoopBackIPs = vars.NewSetOnce(sets.NewString())

	// Subnet is the subnet of the host where the server is located.
	Subnet = vars.SetOnce[string]{}
)
