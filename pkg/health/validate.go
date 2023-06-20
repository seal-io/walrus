package health

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

// MustValidate returns the validation results with the given including list.
func MustValidate(ctx context.Context, includes []string) (string, bool) {
	if len(checkers) == 0 {
		return "no checkers", false
	}

	if len(includes) == 0 {
		return "no include list", false
	}

	var (
		ok = true
		ns = sets.NewString(includes...)
		sb strings.Builder
	)

	for i := range checkers {
		n := checkers[i].Name()

		if !ns.Has(n) {
			continue
		}

		if err := checkers[i].Check(ctx); err != nil {
			ok = false
			_, _ = fmt.Fprintf(&sb, "[-]%s: failed, %v\n", n, err)

			continue
		}

		_, _ = fmt.Fprintf(&sb, "[+]%s: ok\n", n)
	}

	return sb.String(), ok
}

// Validate returns the validation results,
// skips the checker if its name exists in the excluding list.
func Validate(ctx context.Context, excludes ...string) (string, bool) {
	if len(checkers) == 0 {
		return "no checkers", false
	}

	var (
		ok = true
		ns = sets.NewString(excludes...)
		sb strings.Builder
	)

	for i := range checkers {
		n := checkers[i].Name()

		if ns.Has(n) {
			_, _ = fmt.Fprintf(&sb, "[?]%s: excluded\n", n)
			continue
		}

		if err := checkers[i].Check(ctx); err != nil {
			ok = false
			_, _ = fmt.Fprintf(&sb, "[-]%s: failed, %v\n", n, err)

			continue
		}

		_, _ = fmt.Fprintf(&sb, "[+]%s: ok\n", n)
	}

	return sb.String(), ok
}
