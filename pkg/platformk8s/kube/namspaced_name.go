package kube

import "strings"

// ParseNamespacedName parses the given string into {namespace, name},
// e.g. kube-system/coredns.
func ParseNamespacedName(s string) (ns, n string) {
	var ss = strings.SplitN(s, "/", 2)
	if len(ss) == 2 {
		return ss[0], ss[1]
	}
	// use default namespace provided by kubeconfig.
	return "", s
}
