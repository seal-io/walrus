package settings

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/utils/cache"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
)

// Value defines the operations of a built-in setting.
type Value interface {
	// Name returns the name of the setting.
	Name() string

	// Value returns the value of the setting.
	Value(context.Context, model.ClientSet) (string, error)

	// ShouldValue likes Value but without error return,
	// it's good for error-insensitive cases and nice for chain calls.
	ShouldValue(context.Context, model.ClientSet) string

	// ValueJSONUnmarshal unmarshal the setting value into the given holder.
	ValueJSONUnmarshal(context.Context, model.ClientSet, interface{}) error

	// ValueBool returns the bool value of the setting.
	ValueBool(context.Context, model.ClientSet) (bool, error)

	// ShouldValueBool likes ValueBool but without error return,
	// it's good for error-insensitive cases and nice for chain calls.
	ShouldValueBool(context.Context, model.ClientSet) bool

	// ValueInt64 returns the int64 value of the setting.
	ValueInt64(context.Context, model.ClientSet) (int64, error)

	// ShouldValueInt64 likes ValueInt64 but without error return,
	// it's good for error-insensitive cases and nice for chain calls.
	ShouldValueInt64(context.Context, model.ClientSet) int64

	// ValueUint64 returns the uint64 value of the setting.
	ValueUint64(context.Context, model.ClientSet) (uint64, error)

	// ShouldValueUint64 likes ValueUint64 but without error return,
	// it's good for error-insensitive cases and nice for chain calls.
	ShouldValueUint64(context.Context, model.ClientSet) uint64

	// ValueURL returns the url value of the setting.
	ValueURL(context.Context, model.ClientSet) (*url.URL, error)

	// ShouldValueURL likes ValueURL but without error return,
	// it's good for error-insensitive cases and nice for chain calls.
	ShouldValueURL(context.Context, model.ClientSet) *url.URL

	// Set configures the value of the setting,
	// returns true if accept the new value change.
	Set(context.Context, model.ClientSet, interface{}) (bool, error)

	// Cas configures the value of setting with CAS operation.
	Cas(context.Context, model.ClientSet, func(oldVal string) (newVal string, err error)) error
}

var (
	cacher = cache.MustNew(context.Background())
	logger = log.WithName("settings")
)

// value holds the entity implemented the Value interface.
type value struct {
	refer  model.Setting
	modify modifier
}

func (v value) Name() string {
	return v.refer.Name
}

// Value implements the Value interface.
func (v value) Value(ctx context.Context, client model.ClientSet) (string, error) {
	var cachedValue, err = cacher.Get(v.refer.Name)
	if err == nil {
		return string(cachedValue), nil
	}
	dbValue, err := client.Settings().Query().
		Select(setting.FieldValue).
		Where(setting.Name(v.refer.Name)).
		Only(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting %s: %w",
			v.refer.Name, err)
	}
	err = cacher.Set(v.refer.Name, []byte(dbValue.Value))
	if err != nil {
		logger.Warnf("error caching %s: %v",
			v.refer.Name, err)
	}
	return dbValue.Value, nil
}

// ShouldValue implements the Value interface.
func (v value) ShouldValue(ctx context.Context, client model.ClientSet) string {
	var r, _ = v.Value(ctx, client)
	return r
}

// ValueJSONUnmarshal implements the Value interface.
func (v value) ValueJSONUnmarshal(ctx context.Context, client model.ClientSet, holder interface{}) error {
	val, err := v.Value(ctx, client)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), holder)
	if err != nil {
		return fmt.Errorf("error unmarshalling %s: %w",
			v.refer.Name, err)
	}
	return nil
}

// ValueBool implements the Value interface.
func (v value) ValueBool(ctx context.Context, client model.ClientSet) (bool, error) {
	val, err := v.Value(ctx, client)
	if err != nil {
		return false, err
	}
	if val == "" {
		return false, nil
	}
	r, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("error parsing %s: %w",
			v.refer.Name, err)
	}
	return r, nil
}

// ShouldValueBool implements the Value interface.
func (v value) ShouldValueBool(ctx context.Context, client model.ClientSet) bool {
	var r, _ = v.ValueBool(ctx, client)
	return r
}

// ValueInt64 implements the Value interface.
func (v value) ValueInt64(ctx context.Context, client model.ClientSet) (int64, error) {
	val, err := v.Value(ctx, client)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, err
	}
	r, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s: %w",
			v.refer.Name, err)
	}
	return r, nil
}

// ShouldValueInt64 implements the Value interface.
func (v value) ShouldValueInt64(ctx context.Context, client model.ClientSet) int64 {
	var r, _ = v.ValueInt64(ctx, client)
	return r
}

// ValueUint64 implements the Value interface.
func (v value) ValueUint64(ctx context.Context, client model.ClientSet) (uint64, error) {
	val, err := v.Value(ctx, client)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	r, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s: %w",
			v.refer.Name, err)
	}
	return r, nil
}

// ShouldValueUint64 implements the Value interface.
func (v value) ShouldValueUint64(ctx context.Context, client model.ClientSet) uint64 {
	var r, _ = v.ValueUint64(ctx, client)
	return r
}

// ValueURL implements the Value interface.
func (v value) ValueURL(ctx context.Context, client model.ClientSet) (*url.URL, error) {
	val, err := v.Value(ctx, client)
	if err != nil {
		return nil, err
	}
	if val == "" {
		return nil, fmt.Errorf("invalid %s URL: blank",
			v.refer.Name)
	}
	r, err := url.Parse(val)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w",
			v.refer.Name, err)
	}
	return r, nil
}

// ShouldValueURL implements the Value interface.
func (v value) ShouldValueURL(ctx context.Context, client model.ClientSet) *url.URL {
	var r, _ = v.ValueURL(ctx, client)
	return r
}

// Set implements the Value interface.
func (v value) Set(ctx context.Context, client model.ClientSet, newValueRaw interface{}) (bool, error) {
	var oldVal, err = v.Value(ctx, client)
	if err != nil {
		return false, err
	}

	newVal, ok := newValueRaw.(string)
	if !ok {
		b, err := json.Marshal(newValueRaw)
		if err != nil {
			return false, err
		}
		newVal = string(b)
	}

	if oldVal == newVal {
		// nothing to do if same as previous.
		return false, nil
	}

	err = v.modify(ctx, client, v.refer.Name, oldVal, newVal)
	if err != nil {
		return false, err
	}
	err = cacher.Delete(v.refer.Name)
	if err != nil {
		logger.Warnf("error discaching %s: %v",
			v.refer.Name, err)
	}
	return true, nil
}

// Cas implements the Value interface.
func (v value) Cas(ctx context.Context, client model.ClientSet, op func(oldVal string) (newVal string, err error)) error {
	if op == nil {
		return nil
	}
	return client.WithTx(ctx, func(tx *model.Tx) error {
		dbValue, err := tx.Settings().Query().
			Select(setting.FieldValue).
			Where(setting.Name(v.refer.Name)).
			ForUpdate().
			Only(ctx)
		if err != nil {
			return err
		}
		var oldVal = dbValue.Value
		newVal, err := op(oldVal)
		if err != nil {
			return err
		}
		err = v.modify(ctx, tx, v.refer.Name, oldVal, newVal)
		if err != nil {
			return err
		}
		err = cacher.Set(v.refer.Name, []byte(newVal))
		if err != nil {
			logger.Warnf("error caching %s: %v",
				v.refer.Name, err)
		}
		return nil
	})
}
