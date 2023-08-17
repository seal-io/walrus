package validation

import (
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-module/carbon"
	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/seal-io/walrus/utils/errorx"
)

const (
	maxDurationPerYear   = time.Hour * 24 * carbon.DaysPerLeapYear
	maxDurationPerDecade = maxDurationPerYear * carbon.YearsPerDecade
)

func IsDNSLabel(name string) error {
	if len(name) == 0 {
		return errorx.New("name must be non-empty")
	}

	if errs := validation.IsDNS1123Label(name); len(errs) != 0 {
		errStr := strings.Join(errs, ",")
		return errorx.Errorf("name format must conform to DNS Label Names, %s", errStr)
	}

	return nil
}

func TimeRange(startTime, endTime time.Time) error {
	if startTime.IsZero() {
		return errorx.New("invalid start time: blank")
	}

	if endTime.IsZero() {
		return errorx.New("invalid end time: blank")
	}

	if endTime.Before(startTime) {
		return errorx.New("invalid time range: end time is early than start time")
	}

	if startTime.Location().String() != endTime.Location().String() {
		return errorx.New(
			"invalid time range: start time and end time are in different time zones",
		)
	}

	return nil
}

func TimeRangeWithinYear(startTime, endTime time.Time) error {
	if err := TimeRange(startTime, endTime); err != nil {
		return err
	}

	if endTime.Sub(startTime) > maxDurationPerYear {
		return errorx.New(
			"invalid time range: start time and end time must be within a year",
		)
	}

	return nil
}

func TimeRangeWithinDecade(startTime, endTime time.Time) error {
	if err := TimeRange(startTime, endTime); err != nil {
		return err
	}

	if endTime.Sub(startTime) > maxDurationPerDecade {
		return errorx.New(
			"invalid time range: start time and end time must be within decade",
		)
	}

	return nil
}

func IsValidEndpoint(ep string) error {
	if govalidator.IsHost(ep) || govalidator.IsURL(ep) {
		return nil
	}

	return errorx.Errorf("%s isn't a valid endpoint", ep)
}

func IsValidEndpoints(eps []string) error {
	for _, v := range eps {
		if err := IsValidEndpoint(v); err != nil {
			return err
		}
	}

	return nil
}
