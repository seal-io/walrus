package sqlx

import (
	"fmt"

	"github.com/seal-io/walrus/utils/timex"
)

// DateTruncWithZoneOffsetSQL generate the date trunc sql from step and timezone offset,
// offset is in seconds east of UTC.
func DateTruncWithZoneOffsetSQL(field, step string, offset int) (string, error) {
	if step == "" {
		return "", fmt.Errorf("invalid step: blank")
	}

	timezone := timex.TimezoneInPosix(offset)

	switch step {
	case timex.Day, timex.Week, timex.Month, timex.Quarter, timex.Year:
		return fmt.Sprintf(`date_trunc('%s', (%s AT TIME ZONE '%s'))`, step, field, timezone), nil
	default:
		return "", fmt.Errorf("invalid step: unsupport %s", step)
	}
}
