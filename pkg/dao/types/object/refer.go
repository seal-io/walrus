package object

import (
	"strconv"
	"strings"
)

type Refer string

// IsNumeric returns true if the given Refer is a numeric value.
func (i Refer) IsNumeric() bool {
	return isNumeric(string(i))
}

func (i Refer) Int64() int64 {
	v, _ := strconv.ParseInt(string(i), 10, 64)
	return v
}

func (i Refer) Int() int {
	return int(i.Int64())
}

// IsID returns true if the given Refer is an ID value.
func (i Refer) IsID() bool {
	return ID(i).Valid()
}

func (i Refer) ID() ID {
	return ID(i)
}

// IsString returns true if the given Refer is a string value.
func (i Refer) IsString() bool {
	return !i.IsNumeric() && !i.IsID()
}

func (i Refer) String() string {
	return string(i)
}

const defaultSeparator = ":"

// IsComposited returns true if the given Refer is a composited value,
// combines with several fields of the object.
// E.g. X:y:z.
func (i Refer) IsComposited(keyLength int) bool {
	if i.IsNumeric() {
		return false
	}

	return i.matchKeyLength(keyLength)
}

func (i Refer) matchKeyLength(l int) bool {
	switch l {
	case 0:
		return false
	case 1:
		return true
	}

	return strings.Count(string(i), defaultSeparator)+1 == l
}

// Extract extracts the index segment of the composited key with default separator into an array,
// it should be called after IsNumeric == false or IsComposited == true.
func (i Refer) Extract(idx int) ReferSegment {
	return i.Split(idx + 1).Index(idx)
}

// Split splits the composited key with default separator into an array,
// it should be called after IsNumeric == false or IsComposited == true.
func (i Refer) Split(length int) ReferSegments {
	if length <= 1 {
		return []ReferSegment{ReferSegment(i)}
	}

	ss := strings.SplitN(string(i), defaultSeparator, length)
	if len(ss) != length {
		return []ReferSegment{}
	}

	v := make([]ReferSegment, length)
	for i := range ss {
		v[i] = ReferSegment(ss[i])
	}

	return v
}

type ReferSegment string

func (i ReferSegment) ID() ID {
	return ID(i)
}

func (i ReferSegment) Int64() int64 {
	v, _ := strconv.ParseInt(string(i), 10, 64)
	return v
}

func (i ReferSegment) Int() int {
	return int(i.Int64())
}

func (i ReferSegment) Bool() bool {
	v, _ := strconv.ParseBool(string(i))
	return v
}

func (i ReferSegment) String() string {
	return string(i)
}

type ReferSegments []ReferSegment

// Index gets the segment with the given order.
func (i ReferSegments) Index(idx int) ReferSegment {
	if 0 > idx || idx >= len(i) {
		return ""
	}

	return i[idx]
}
