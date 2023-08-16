package key

import (
	"strings"

	"github.com/seal-io/walrus/utils/strs"
)

// Decode parses the given string into {resource type, resource name}, returns false if not a valid key.
func Decode(s string) (resourceType, name string, ok bool) {
	ss := strings.SplitN(s, "/", 2)

	ok = len(ss) == 2
	if !ok {
		return
	}
	resourceType = ss[0]
	name = ss[1]

	return
}

// Encode constructs the given {resource type, resource name} into a valid key.
func Encode(resourceType, name string) string {
	return strs.Join("/", resourceType, name)
}
