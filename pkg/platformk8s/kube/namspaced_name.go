package kube

import (
	"strings"

	"github.com/seal-io/seal/utils/strs"
)

// ParseNamespacedName parses the given string into {namespace, name},
// e.g. kube-system/coredns.
func ParseNamespacedName(s string) (ns, n string) {
	var ss = strings.SplitN(s, "/", 2)
	if len(ss) == 2 {
		return ss[0], ss[1]
	}
	// Use default namespace provided by kubeconfig.
	return "", s
}

// NamespacedName constructs the given {namespace, name} into one string.
func NamespacedName(ns, n string) string {
	if ns == "" {
		return n
	}
	return strs.Join("/", ns, n)
}
