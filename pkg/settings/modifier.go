package settings

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	imgdistref "github.com/distribution/distribution/reference"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/utils/cron"
)

// modifier defines the stereotype for writing value.
type modifier func(ctx context.Context, client model.ClientSet, name, oldValue, newValue string) error

func modifyWith(validates ...modifyValidator) modifier {
	return func(ctx context.Context, client model.ClientSet, name, oldValue, newValue string) error {
		if len(validates) == 0 {
			validates = append(validates, many)
		}

		for i := range validates {
			ok, err := validates[i](ctx, name, oldValue, newValue)
			if err != nil {
				return runtime.Errorf(http.StatusBadRequest, "invalid setting %q: %w", name, err)
			}

			if !ok {
				return nil
			}
		}

		err := client.Settings().Update().
			SetValue(crypto.String(newValue)).
			Where(setting.Name(name)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("error modify setting %s: %w", name, err)
		}

		return nil
	}
}

// modifyValidator defines the stereotype for validating writing.
type modifyValidator func(ctx context.Context, name, oldVal, newVal string) (bool, error)

// notBlank implements the modifyValidator stereotype,
// which means the value can be modified if not blank.
func notBlank(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	if isBlankValue(newVal) {
		return false, errors.New("blank value")
	}

	return true, nil
}

// many implements the modifyValidator stereotype,
// which means the value can be modified if not the same.
func many(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	return oldVal != newVal, nil
}

// once implements the modifyValidator stereotype,
// which means the value can be modified if blank.
func once(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	if !isBlankValue(oldVal) {
		return false, errors.New("already configured")
	}

	return true, nil
}

// never implements the modifyValidator stereotype,
// which means the value can not be modified.
func never(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	return false, errors.New("cannot modify")
}

// httpUrl implements the modifyValidator stereotype,
// which means the value can be modified if it is an HTTP URL.
// This modifier allows blank new value,
// if not allowed, combine with notBlank.
func httpUrl(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	if newVal == "" {
		return true, nil
	}

	_, err := parseUrl(newVal, httpSchemeUrlOnly)

	return err == nil, err
}

// sockUrl implements the modifyValidator stereotype,
// which means the value can be modified if it is a Socket URL.
// This modifier allows blank new value,
// if not allowed, combine with notBlank.
func sockUrl(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	// Allow blank,
	// combine with notBlank if disallowed.
	if newVal == "" {
		return true, nil
	}

	_, err := parseUrl(newVal, sockSchemeUrlOnly)

	return err == nil, err
}

// anyUrl implements the modifyValidator stereotype,
// which means the value can be modified if it is a URL.
// This modifier allows blank new value,
// if not allowed, combine with notBlank.
func anyUrl(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	// Allow blank,
	// combine with notBlank if disallowed.
	if newVal == "" {
		return true, nil
	}

	_, err := parseUrl(newVal, anySchemeUrl)

	return err == nil, err
}

// cronExpression implements the modifyValidator stereotype,
// which means the value can be modified if it's cron expression.
// This modifier allows blank new value,
// if not allowed, combine with notBlank.
func cronExpression(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	// Allow blank,
	// combine with notBlank if disallowed.
	if newVal == "" {
		return true, nil
	}

	err := cron.ValidateCronExpr(newVal)

	return err == nil, err
}

// containerImageReference implements the modifyValidator stereotype,
// which means the value can be modified if it's container image reference.
// This modifier allows blank new value,
// if not allowed, combine with notBlank.
func containerImageReference(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	// Allow blank,
	// combine with notBlank if disallowed.
	if newVal == "" {
		return true, nil
	}

	_, err := imgdistref.ParseNormalizedNamed(newVal)

	return err == nil, err
}
