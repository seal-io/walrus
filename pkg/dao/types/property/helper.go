package property

import (
	"time"

	"github.com/seal-io/walrus/utils/json"
)

func GetNumber(i Value) (any, bool, error) {
	iv, ok, _ := GetInt64(i)
	if ok {
		return iv, ok, nil
	}

	uiv, ok, _ := GetUint64(i)
	if ok {
		return uiv, ok, nil
	}

	fv, ok, err := GetFloat64(i)
	if ok {
		return fv, ok, nil
	}

	return 0, false, err
}

func GetUint64(i Value) (uint64, bool, error) {
	if i != nil {
		var v uint64
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return 0, false, nil
}

func GetUint32(i Value) (uint32, bool, error) {
	v, ok, err := GetUint64(i)
	if err == nil && ok {
		return uint32(v), true, nil
	}

	return 0, ok, err
}

func GetUint16(i Value) (uint16, bool, error) {
	v, ok, err := GetUint64(i)
	if err == nil && ok {
		return uint16(v), true, nil
	}

	return 0, ok, err
}

func GetUint8(i Value) (uint8, bool, error) {
	v, ok, err := GetUint64(i)
	if err == nil && ok {
		return uint8(v), true, nil
	}

	return 0, ok, err
}

func GetUint(i Value) (uint, bool, error) {
	v, ok, err := GetUint64(i)
	if err == nil && ok {
		return uint(v), true, nil
	}

	return 0, ok, err
}

func GetInt64(i Value) (int64, bool, error) {
	if i != nil {
		var v int64
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return 0, false, nil
}

func GetInt32(i Value) (int32, bool, error) {
	v, ok, err := GetInt64(i)
	if err == nil && ok {
		return int32(v), true, nil
	}

	return 0, ok, err
}

func GetInt16(i Value) (int16, bool, error) {
	v, ok, err := GetInt64(i)
	if err == nil && ok {
		return int16(v), true, nil
	}

	return 0, ok, err
}

func GetInt8(i Value) (int8, bool, error) {
	v, ok, err := GetInt64(i)
	if err == nil && ok {
		return int8(v), true, nil
	}

	return 0, ok, err
}

func GetInt(i Value) (int, bool, error) {
	v, ok, err := GetInt64(i)
	if err == nil && ok {
		return int(v), true, nil
	}

	return 0, ok, err
}

func GetFloat64(i Value) (float64, bool, error) {
	if i != nil {
		var v float64
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return 0, false, nil
}

func GetFloat32(i Value) (float32, bool, error) {
	v, ok, err := GetFloat64(i)
	if err == nil && ok {
		return float32(v), true, nil
	}

	return 0, ok, err
}

func GetDuration(i Value) (time.Duration, bool, error) {
	if i != nil {
		var v time.Duration
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return 0, false, nil
}

func GetBool(i Value) (bool, bool, error) {
	if i != nil {
		var v bool
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return false, false, nil
}

func GetString(i Value) (string, bool, error) {
	if i != nil {
		var v string
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return "", false, nil
}

// GetSlice returns the underlay value as a slice with the given generic type,
// if not found or parse error, returns false.
func GetSlice[T any](i Value) ([]T, bool, error) {
	if i != nil {
		var v []T
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return nil, false, nil
}

// GetMap returns the underlay value as a string map with the given generic type,
// if not found or parse error, returns false.
func GetMap[T any](i Value) (map[string]T, bool, error) {
	if i != nil {
		var v map[string]T
		err := json.Unmarshal(i, &v)

		return v, err == nil, err
	}

	return nil, false, nil
}

// GetObject returns the underlay value as a T object with the given generic type,
// if not found or parse error, returns false.
func GetObject[T any](i Value) (T, bool, error) {
	var v T

	if i != nil {
		err := json.Unmarshal(i, &v)
		return v, err == nil, err
	}

	return v, false, nil
}

// GetAny returns the underlay value as the given generic type,
// if not found or parse error, returns false.
func GetAny[T any](i Value) (T, bool, error) {
	var v T

	if i != nil {
		err := json.Unmarshal(i, &v)
		return v, err == nil, err
	}

	return v, false, nil
}
