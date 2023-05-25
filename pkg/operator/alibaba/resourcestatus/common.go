package resourcestatus

import (
	"errors"
	"fmt"
)

const schemeHttps = "HTTPS"

var errNotFound = errors.New("not found")

func toReqIds(id string) string {
	return fmt.Sprintf(`["%s"]`, id)
}
