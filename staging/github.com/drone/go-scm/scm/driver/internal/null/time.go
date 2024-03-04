// Copyright (c) 2014, Greg Roseberry
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package null

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Time is a nullable time.Time.
// JSON marshals to the zero value for time.Time if null.
// Considered to be null to SQL if zero.
type Time struct {
	Time  time.Time
	Valid bool
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case string:
		err = t.Time.UnmarshalJSON(data)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// IsZero returns true for null or zero Times, for
// potential future omitempty support.
func (t Time) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// Ptr returns a pointer to this Time's value,
// or a nil pointer if this Time is zero.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// NewTime creates a new Time.
func NewTime(t time.Time, valid bool) Time {
	return Time{
		Time:  t,
		Valid: valid,
	}
}

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time {
	return NewTime(t, !t.IsZero())
}