package crypto

import (
	"database/sql/driver"
	"errors"

	"github.com/zclconf/go-cty/cty"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/strs"
)

const sensitive = "<sensitive>"

// String shows the secret value in string but stores in byte array.
type String string

// String implements fmt.Stringer.
func (i String) String() string {
	return sensitive
}

// Value implements driver.Valuer.
func (i String) Value() (driver.Value, error) {
	v := string(i)
	enc := EncryptorConfig.Get()

	return enc.Encrypt(strs.ToBytes(&v), nil)
}

// Scan implements sql.Scanner.
func (i *String) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		enc := EncryptorConfig.Get()

		p, err := enc.Decrypt(v, nil)
		if err != nil {
			return err
		}
		*i = String(strs.FromBytes(&p))

		return nil
	}

	return errors.New("not a valid crypto string")
}

// Bytes encrypts/decrypts the byte array.
type Bytes []byte

// String implements fmt.Stringer.
func (i Bytes) String() string {
	return sensitive
}

// Value implements driver.Valuer.
func (i Bytes) Value() (driver.Value, error) {
	v := i
	enc := EncryptorConfig.Get()

	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Bytes) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		enc := EncryptorConfig.Get()

		p, err := enc.Decrypt(v, nil)
		if err != nil {
			return err
		}
		*i = p

		return nil
	}

	return errors.New("not a valid crypto bytes")
}

// Map shows the secret value in map but stores in byte array.
type Map[K comparable, V any] map[K]V

// String implements fmt.Stringer.
func (i Map[K, V]) String() string {
	return sensitive
}

// Value implements driver.Valuer.
func (i Map[K, V]) Value() (driver.Value, error) {
	v, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	enc := EncryptorConfig.Get()

	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Map[K, V]) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		enc := EncryptorConfig.Get()

		p, err := enc.Decrypt(v, nil)
		if err != nil {
			return err
		}

		return json.Unmarshal(p, i)
	}

	return errors.New("not a valid crypto map")
}

// Slice shows the secret value in slice but stores in byte array.
type Slice[T any] []T

// String implements fmt.Stringer.
func (i Slice[T]) String() string {
	return sensitive
}

// Value implements driver.Valuer.
func (i Slice[T]) Value() (driver.Value, error) {
	v, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	enc := EncryptorConfig.Get()

	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Slice[T]) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		enc := EncryptorConfig.Get()

		p, err := enc.Decrypt(v, nil)
		if err != nil {
			return err
		}

		return json.Unmarshal(p, i)
	}

	return errors.New("not a valid crypto slice")
}

// Property wraps the property.Property,
// and erases the invisible value during output.
type Property struct {
	property.Property `json:",inline"`

	// Visible indicates to show the value of this property or not.
	Visible bool `json:"visible"`
}

// String implements fmt.Stringer.
func (i Property) String() string {
	if !i.Visible {
		return sensitive
	}

	return string(i.Value)
}

// MarshalJSON implements json.Marshaler,
// impacts the response message.
func (i Property) MarshalJSON() ([]byte, error) {
	type Alias Property

	ia := (Alias)(i)
	if !i.Visible {
		ia.Value = nil
	}

	return json.Marshal(ia)
}

// Properties holds the secret Property collection in map,
// the key of map is the name of Property,
// stores into byte array.
type Properties map[string]Property

// Value implements driver.Valuer.
func (i Properties) Value() (driver.Value, error) {
	type Alias Property

	ia := make(map[string]Alias, len(i))
	for k := range i {
		ia[k] = (Alias)(i[k])
	}

	v, err := json.Marshal(ia)
	if err != nil {
		return nil, err
	}
	enc := EncryptorConfig.Get()

	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Properties) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		enc := EncryptorConfig.Get()

		p, err := enc.Decrypt(v, nil)
		if err != nil {
			return err
		}

		return json.Unmarshal(p, i)
	}

	return errors.New("not a valid crypto properties")
}

// Cty returns the type and value of this property.
func (i Properties) Cty() (cty.Type, cty.Value, error) {
	var (
		ot = make(map[string]cty.Type, len(i))
		ov = make(map[string]cty.Value, len(i))
	)

	for x := range i {
		t, v, err := i[x].Cty()
		if err != nil {
			return cty.NilType, cty.NilVal, err
		}
		ot[x] = t
		ov[x] = v
	}

	return cty.Object(ot), cty.ObjectVal(ov), nil
}

// Values returns a map stores the underlay value.
func (i Properties) Values() property.Values {
	m := make(property.Values, len(i))
	for x := range i {
		m[x] = i[x].GetValue()
	}

	return m
}

// TypedValues returns a map stores the typed value.
func (i Properties) TypedValues() (m map[string]any, err error) {
	m = make(map[string]any, len(i))

	for x := range i {
		typ := i[x].GetType()

		switch {
		case typ == cty.Number:
			m[x], _, err = i[x].GetNumber()
		case typ == cty.Bool:
			m[x], _, err = i[x].GetBool()
		case typ == cty.String:
			m[x], _, err = i[x].GetString()
		case typ.IsListType() || typ.IsTupleType() || typ.IsSetType():
			m[x], _, err = property.GetSlice[any](i[x])
		case typ.IsMapType() || typ.IsObjectType():
			m[x], _, err = property.GetMap[any](i[x])
		default:
			m[x], _, err = property.GetAny[any](i[x])
		}

		if err != nil {
			return
		}
	}

	return
}
