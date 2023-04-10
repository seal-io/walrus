package crypto

import (
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/types/property"
)

// Uint64Property wraps uint64 value into a property.
func Uint64Property(v uint64) Property {
	return Property{
		Property: property.Uint64Property(v),
	}
}

// Uint32Property wraps uint32 value into a property.
func Uint32Property(v uint32) Property {
	return Property{
		Property: property.Uint32Property(v),
	}
}

// Uint16Property wraps uint16 value into a property.
func Uint16Property(v uint16) Property {
	return Property{
		Property: property.Uint16Property(v),
	}
}

// Uint8Property wraps uint8 value into a property.
func Uint8Property(v uint8) Property {
	return Property{
		Property: property.Uint8Property(v),
	}
}

// UintProperty wraps uint value into a property.
func UintProperty(v uint) Property {
	return Property{
		Property: property.UintProperty(v),
	}
}

// Int64Property wraps int64 value into a property.
func Int64Property(v int64) Property {
	return Property{
		Property: property.Int64Property(v),
	}
}

// Int32Property wraps int32 value into a property.
func Int32Property(v int32) Property {
	return Property{
		Property: property.Int32Property(v),
	}
}

// Int16Property wraps int16 value into a property.
func Int16Property(v int16) Property {
	return Property{
		Property: property.Int16Property(v),
	}
}

// Int8Property wraps int8 value into a property.
func Int8Property(v int8) Property {
	return Property{
		Property: property.Int8Property(v),
	}
}

// IntProperty wraps int value into a property.
func IntProperty(v int) Property {
	return Property{
		Property: property.IntProperty(v),
	}
}

// Float64Property wraps float64 value into a property.
func Float64Property(v float64) Property {
	return Property{
		Property: property.Float64Property(v),
	}
}

// Float32Property wraps float32 value into a property.
func Float32Property(v float32) Property {
	return Property{
		Property: property.Float32Property(v),
	}
}

// DurationProperty wraps time.Duration value into a property.
func DurationProperty(v time.Duration) Property {
	return Property{
		Property: property.DurationProperty(v),
	}
}

// BoolProperty wraps bool value into a property.
func BoolProperty(v bool) Property {
	return Property{
		Property: property.BoolProperty(v),
	}
}

// StringProperty wraps string value into a property.
func StringProperty(v string) Property {
	return Property{
		Property: property.StringProperty(v),
	}
}

// SliceProperty wraps slice value into a property.
func SliceProperty[T any](v []T) Property {
	return Property{
		Property: property.SliceProperty[T](v),
	}
}

// SetProperty wraps set value into a property.
func SetProperty[T comparable](v sets.Set[T]) Property {
	return Property{
		Property: property.SetProperty[T](v),
	}
}

// MapProperty wraps map value into a property.
func MapProperty[T any](v map[string]T) Property {
	return Property{
		Property: property.MapProperty[T](v),
	}
}

// ObjectProperty wraps object value into a property.
func ObjectProperty[T any](v T) Property {
	return Property{
		Property: property.ObjectProperty[T](v),
	}
}

// AnyProperty wraps any value into a property.
func AnyProperty(v any) Property {
	return Property{
		Property: property.AnyProperty(v),
	}
}
