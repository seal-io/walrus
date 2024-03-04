package stringx

import (
	"strconv"
)

// integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type integer interface {
	signed | unsigned
}

// signed is a constraint that permits any signed integer type.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// unsigned is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// FromInt is similar to strconv.Itoa,
// but it accepts any integer type.
func FromInt[T integer](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

// ToInt is similar to strconv.Atoi,
// but it returns any integer type.
func ToInt[T integer](s string) (T, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return T(i), err
}
