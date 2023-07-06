package validation

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-module/carbon"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	maxDurationPerYear   = time.Hour * 24 * carbon.DaysPerLeapYear
	maxDurationPerDecade = maxDurationPerYear * carbon.YearsPerDecade
)

func IsDNSLabel(name string) error {
	if len(name) == 0 {
		return errors.New("name must be non-empty")
	}

	if errs := validation.IsDNS1123Label(name); len(errs) != 0 {
		errStr := strings.Join(errs, ",")
		return fmt.Errorf("name format must conform to DNS Label Names, %s", errStr)
	}

	return nil
}

func TimeRange(startTime, endTime time.Time) error {
	if startTime.IsZero() {
		return errors.New("invalid start time: blank")
	}

	if endTime.IsZero() {
		return errors.New("invalid end time: blank")
	}

	if endTime.Before(startTime) {
		return errors.New("invalid time range: end time is early than start time")
	}

	if startTime.Location().String() != endTime.Location().String() {
		return errors.New("invalid time range: start time and end time are in different time zones")
	}

	return nil
}

func TimeRangeWithinYear(startTime, endTime time.Time) error {
	if err := TimeRange(startTime, endTime); err != nil {
		return err
	}

	if endTime.Sub(startTime) > maxDurationPerYear {
		return fmt.Errorf("invalid time range: start time and end time must be within a year")
	}

	return nil
}

func TimeRangeWithinDecade(startTime, endTime time.Time) error {
	if err := TimeRange(startTime, endTime); err != nil {
		return err
	}

	if endTime.Sub(startTime) > maxDurationPerDecade {
		return fmt.Errorf("invalid time range: start time and end time must be within decade")
	}

	return nil
}

func IsValidEndpoint(ep string) error {
	if govalidator.IsHost(ep) || govalidator.IsURL(ep) {
		return nil
	}

	return fmt.Errorf("%s isn't a valid endpoint", ep)
}

func IsValidEndpoints(eps []string) error {
	for _, v := range eps {
		if err := IsValidEndpoint(v); err != nil {
			return err
		}
	}

	return nil
}
