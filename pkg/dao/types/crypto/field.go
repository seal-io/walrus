package crypto

import (
	"database/sql/driver"
	"errors"

	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/strs"
)

// String shows the secret value in string but stores in byte array.
type String string

// String implements fmt.Stringer.
func (i String) String() string {
	return "<sensitive>"
}

// Value implements driver.Valuer.
func (i String) Value() (driver.Value, error) {
	var v = string(i)
	var enc = EncryptorConfig.Get()
	return enc.Encrypt(strs.ToBytes(&v), nil)
}

// Scan implements sql.Scanner.
func (i *String) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		var enc = EncryptorConfig.Get()
		var p, err = enc.Decrypt(v, nil)
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
	return "<sensitive>"
}

// Value implements driver.Valuer.
func (i Bytes) Value() (driver.Value, error) {
	var v = i
	var enc = EncryptorConfig.Get()
	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Bytes) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		var enc = EncryptorConfig.Get()
		var p, err = enc.Decrypt(v, nil)
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
	return "<sensitive>"
}

// Value implements driver.Valuer.
func (i Map[K, V]) Value() (driver.Value, error) {
	var v, err = json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var enc = EncryptorConfig.Get()
	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Map[K, V]) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		var enc = EncryptorConfig.Get()
		var p, err = enc.Decrypt(v, nil)
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
	return "<sensitive>"
}

// Value implements driver.Valuer.
func (i Slice[T]) Value() (driver.Value, error) {
	var v, err = json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var enc = EncryptorConfig.Get()
	return enc.Encrypt(v, nil)
}

// Scan implements sql.Scanner.
func (i *Slice[T]) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		var enc = EncryptorConfig.Get()
		var p, err = enc.Decrypt(v, nil)
		if err != nil {
			return err
		}
		return json.Unmarshal(p, i)
	}
	return errors.New("not a valid crypto slice")
}
