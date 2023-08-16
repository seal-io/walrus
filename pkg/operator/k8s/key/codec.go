package key

import (
	"strings"

	"github.com/seal-io/walrus/utils/strs"
)

// Decode parses the given string into {pod namespace, pod name, container type, container name},
// returns false if not a valid key, e.g. default/coredns-64897985d-6x2jm/container/coredns.
// Valid container types have `initContainer`, `ephemeralContainer`, `container`.
func Decode(s string) (podNamespace, podName, containerType, containerName string, ok bool) {
	ss := strings.SplitN(s, "/", 4)

	ok = len(ss) == 4
	if !ok {
		return
	}
	podNamespace = ss[0]
	podName = ss[1]
	containerType = ss[2]
	containerName = ss[3]

	return
}

// Encode constructs the given {pod namespace, pod name, container type, container name} into a valid key.
func Encode(podNamespace, podName, containerType, containerName string) string {
	return strs.Join("/", podNamespace, podName, containerType, containerName)
}
