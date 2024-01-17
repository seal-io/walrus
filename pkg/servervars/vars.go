package servervars

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/vars"
)

// NonLoopBackIPs is a set of non-loopback local IPs.
var NonLoopBackIPs = vars.NewSetOnce(sets.New[string]())
