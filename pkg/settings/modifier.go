package settings

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/setting"
)

// modifier defines the stereotype for writing value.
type modifier func(ctx context.Context, client model.ClientSet, name, oldValue, newValue string) error

func modifyWith(validates ...modifyValidator) modifier {
	return func(ctx context.Context, client model.ClientSet, name, oldValue, newValue string) error {
		if len(validates) == 0 {
			validates = append(validates, many)
		}
		for i := range validates {
			var ok, err = validates[i](ctx, name, oldValue, newValue)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
		}
		var err = client.Settings().Update().
			SetValue(newValue).
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

// many implements the modifyValidator stereotype,
// which means the value can be modified if not the same.
func many(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	return oldVal == newVal, nil
}

// once implements the modifyValidator stereotype,
// which means the value can be modified if blank.
func once(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	if oldVal != "" && oldVal != "{}" && oldVal != "[]" {
		return false, fmt.Errorf("setting %s has been configured", name)
	}
	return true, nil
}

// httpUrl implements the modifyValidator stereotype,
// which means the value can be modified if it is an HTTP URL.
func httpUrl(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	var _, err = parseUrl(newVal, httpSchemeUrlOnly)
	if err != nil {
		return false, err
	}
	return true, nil
}

// sockUrl implements the modifyValidator stereotype,
// which means the value can be modified if it is a Socket URL.
func sockUrl(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	var _, err = parseUrl(newVal, sockSchemeUrlOnly)
	if err != nil {
		return false, err
	}
	return true, nil
}

// anyUrl implements the modifyValidator stereotype,
// which means the value can be modified if it is an URL.
func anyUrl(ctx context.Context, name, oldVal, newVal string) (bool, error) {
	var _, err = parseUrl(newVal, anySchemeUrl)
	if err != nil {
		return false, err
	}
	return true, nil
}
