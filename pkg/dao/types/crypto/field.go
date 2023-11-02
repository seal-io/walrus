package crypto

import (
	"database/sql/driver"
	"errors"

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
	// Value specifies the value of this property.
	Value property.Value `json:"value,omitempty"`
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
