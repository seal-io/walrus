package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-module/carbon"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	qnameCharFmt        string = "[A-Za-z0-9]"
	qnameExtCharFmt     string = "[-A-Za-z0-9_.]"
	qualifiedNameFmt    string = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
	qualifiedNameErrMsg string = "a qualified name must consist of alphanumeric characters, '-'," +
		" '_' or '.', and must start and end with an alphanumeric character(e.g. MyName, my.name or 123-abc)"
)
const qualifiedNameMaxLength int = 60

const (
	maxDurationPerYear   = time.Hour * 24 * carbon.DaysPerLeapYear
	maxDurationPerDecade = maxDurationPerYear * carbon.YearsPerDecade
)

var qualifiedNameRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")

func IsQualifiedName(name string) error {
	if len(name) == 0 {
		return errors.New("name must be non-empty")
	} else if len(name) > qualifiedNameMaxLength {
		return fmt.Errorf("name must be no more than %d characters", qualifiedNameMaxLength)
	}

	if !qualifiedNameRegexp.MatchString(name) {
		return fmt.Errorf("%s, regex used for validation is '%s'", qualifiedNameErrMsg, qualifiedNameFmt)
	}

	return nil
}

func IsDNSSubdomainName(name string) error {
	if len(name) == 0 {
		return errors.New("name must be non-empty")
	} else if len(name) > qualifiedNameMaxLength {
		return fmt.Errorf("name must be no more than %d characters", qualifiedNameMaxLength)
	}

	if errs := validation.IsDNS1123Subdomain(name); len(errs) != 0 {
		errStr := strings.Join(errs, ",")
		return fmt.Errorf("name format must conform to DNS Subdomain Names, %s", errStr)
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
