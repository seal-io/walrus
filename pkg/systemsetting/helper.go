package systemsetting

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/distribution/reference"
	"github.com/robfig/cron/v3"
	"github.com/seal-io/utils/osx"
	"github.com/seal-io/utils/stringx"
	"k8s.io/apimachinery/pkg/util/validation"
)

// Initializer defines the stereotype for initializing the default value.
type Initializer func(name string) (value string)

// InitializeFromSpecifiedEnv searches the variable of the given environment name,
// returns the default value if not found the specified environment.
func InitializeFromSpecifiedEnv(envName string, defValue ...string) Initializer {
	return func(_ string) string {
		return osx.ExpandEnv(envName, defValue...)
	}
}

// InitializeFromEnv searches the variable of the `WALRUS_SETTING_${UpperSnakeCase_of_SettingName}`,
// returns the default value if not found the environment.
func InitializeFromEnv(defValue ...string) Initializer {
	return func(name string) string {
		envName := "WALRUS_SETTING_" + stringx.UnderscoreUpper(name)
		return osx.ExpandEnv(envName, defValue...)
	}
}

// InitializeFrom initializes with the given value.
func InitializeFrom(value string) Initializer {
	return func(_ string) string {
		return value
	}
}

// Admission defines the stereotype for validating writing.
type Admission func(ctx context.Context, oldVal, newVal string) (err error)

// ErrAdmissionSkipped is the error returned by Admission to indicate the validation is skipped.
var ErrAdmissionSkipped = errors.New("admission skipped")

// AdmitWith combines multiple Admission into one,
// if no Admission given, it will be Allow.
//
// If one Admission of the given Admissions returns an error,
// the error will be returned immediately.
//
// If one Admission of the given Admissions returns an ErrAdmissionSkipped,
// the function will return true immediately.
func AdmitWith(admits ...Admission) Admission {
	if len(admits) == 0 {
		admits = []Admission{Allow}
	}
	return func(ctx context.Context, oldValue, newValue string) error {
		var err error
		for i := range admits {
			err = admits[i](ctx, oldValue, newValue)
			if err != nil {
				break
			}
		}
		if errors.Is(err, ErrAdmissionSkipped) {
			return nil
		}
		return err
	}
}

// Allow implements the Admission stereotype,
// which means the value can be modified.
func Allow(ctx context.Context, oldVal, newVal string) error {
	return nil
}

// Disallow implements the Admission stereotype,
// which means the value can not be modified.
func Disallow(ctx context.Context, oldVal, newVal string) error {
	return errors.New("cannot modify")
}

// AllowBlank implements the Admission stereotype,
// which means the value can be modified if blank.
//
// AllowBlank always combines with other Admission,
// if allow input anything, please use Allow instead.
func AllowBlank(ctx context.Context, oldVal, newVal string) error {
	if !isBlank(newVal) {
		return nil
	}
	return ErrAdmissionSkipped
}

// DisallowBlank implements the Admission stereotype,
// which means the value can be modified if not blank.
//
// DisallowBlank always combines with other Admission,
// if disallow input anything, please use Disallow instead.
func DisallowBlank(ctx context.Context, oldVal, newVal string) error {
	if !isBlank(newVal) {
		return nil
	}
	return errors.New("blank value")
}

// AllowOnceConfigure implements the Admission stereotype,
// which means the value can be modified if blank.
func AllowOnceConfigure(ctx context.Context, oldVal, newVal string) error {
	if isBlank(oldVal) {
		return nil
	}
	return errors.New("already configured")
}

// AllowBoolean implements the Admission stereotype,
// which means the value can be modified if it's boolean.
func AllowBoolean(ctx context.Context, oldVal, newVal string) error {
	_, err := strconv.ParseBool(newVal)
	return err
}

// AllowInt64 implements the Admission stereotype,
// which means the value can be modified if it's int64.
func AllowInt64(ctx context.Context, oldVal, newVal string) error {
	_, err := strconv.ParseInt(newVal, 10, 64)
	return err
}

// AllowUint64 implements the Admission stereotype,
// which means the value can be modified if it's uint64.
func AllowUint64(ctx context.Context, oldVal, newVal string) error {
	_, err := strconv.ParseUint(newVal, 10, 64)
	return err
}

// AllowFloat64 implements the Admission stereotype,
// which means the value can be modified if it's float64.
func AllowFloat64(ctx context.Context, oldVal, newVal string) error {
	_, err := strconv.ParseFloat(newVal, 64)
	return err
}

// AllowHttpUrl implements the Admission stereotype,
// which means the value can be modified if it is an HTTP URL.
// This Admission allows blank new value,
// if not allowed, combine with DisallowBlank.
func AllowHttpUrl(ctx context.Context, oldVal, newVal string) error {
	return checkUrl(newVal, httpSchemeUrlOnly)
}

// AllowSockUrl implements the Admission stereotype,
// which means the value can be modified if it is a Socket URL.
// This Admission allows blank new value,
// if not allowed, combine with DisallowBlank.
func AllowSockUrl(ctx context.Context, oldVal, newVal string) error {
	return checkUrl(newVal, sockSchemeUrlOnly)
}

// AllowUrl implements the Admission stereotype,
// which means the value can be modified if it is a URL.
// This Admission allows blank new value,
// if not allowed, combine with DisallowBlank.
func AllowUrl(ctx context.Context, oldVal, newVal string) error {
	return checkUrl(newVal, anySchemeUrl)
}

// AllowCronExpression implements the Admission stereotype,
// which means the value can be modified if it's cron expression.
// This Admission allows blank new value,
// if not allowed, combine with DisallowBlank.
func AllowCronExpression(ctx context.Context, oldVal, newVal string) error {
	_, err := cron.ParseStandard(newVal)
	return err
}

// AllowCronExpressionAtLeast implements the Admission stereotype,
// which means the value can be modified if it's cron expression and at least the given duration.
// This Admission allows blank new value,
// if not allowed, combine with DisallowBlank.
func AllowCronExpressionAtLeast(d time.Duration) Admission {
	return func(ctx context.Context, oldVal, newVal string) error {
		expr, err := cron.ParseStandard(newVal)
		if err != nil {
			return err
		}
		next := expr.Next(time.Now())
		duration := expr.Next(next).Sub(next)
		if duration < d {
			return fmt.Errorf("cron expression %q is too short, at least %v", newVal, d)
		}
		return nil
	}
}

// AllowContainerRegistry implements the Admission stereotype,
// which means the value can be modified if it's container registry.
func AllowContainerRegistry(ctx context.Context, oldVal, newVal string) error {
	u, err := url.Parse(newVal)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	if u.Host == "" {
		return errors.New("empty host")
	}
	if msg := validation.IsQualifiedName(newVal); len(msg) != 0 {
		return fmt.Errorf("invalid: %v", msg)
	}
	return nil
}

// AllowContainerImageReference implements the Admission stereotype,
// which means the value can be modified if it's container image reference.
func AllowContainerImageReference(ctx context.Context, oldVal, newVal string) error {
	_, err := reference.ParseNormalizedNamed(newVal)
	return err
}

func isBlank(s string) bool {
	return slices.Contains([]string{"", "{}", "[]"}, strings.TrimSpace(s))
}

type urlChecker func(url.URL) error

func anySchemeUrl(_ url.URL) error {
	return nil
}

func httpSchemeUrlOnly(u url.URL) error {
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		return nil
	}
	return fmt.Errorf("invalid schema: %q, allowed: http, https", u.Scheme)
}

func sockSchemeUrlOnly(u url.URL) error {
	switch strings.ToLower(u.Scheme) {
	case "socks4", "socks5":
		return nil
	}
	return fmt.Errorf("invalid schema: %q, allowed socks4, socks5", u.Scheme)
}

func checkUrl(str string, check urlChecker) error {
	v, err := url.Parse(str)
	if err != nil {
		return fmt.Errorf("%s is illegal URL format: %w", str, err)
	}

	port := v.Port()
	if port != "" {
		p, err := strconv.ParseUint(port, 10, 32)
		if err != nil {
			return fmt.Errorf("parsing given port: %w", err)
		}

		if p > 65535 {
			return fmt.Errorf("given port %d: exceeded upper limit", p)
		}
	}

	return check(*v)
}
