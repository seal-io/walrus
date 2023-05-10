package oid

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"
)

func New(id uint64) ID {
	return ID(strconv.FormatUint(id, 10))
}

// ID shows the primary key in string but stores in big integer,
// also be good at catching composited primary keys.
type ID string

// String implements fmt.Stringer.
func (i ID) String() string {
	return string(i)
}

// Value implements driver.Valuer.
func (i ID) Value() (driver.Value, error) {
	r, err := strconv.ParseInt(string(i), 10, 64)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Scan implements sql.Scanner.
func (i *ID) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case int64:
		*i = ID(strconv.FormatInt(v, 10))
		return nil
	}
	return errors.New("not a valid ID")
}

// Valid returns true if the given ID is a naive ID or composited ID,
// giving zero key length means the ID must be a naive ID.
func (i ID) Valid(keyLength int) bool {
	if i.IsNaive() {
		return true
	}
	return i.matchKeyLength(keyLength)
}

// IsNaive returns true if the given ID is a naive ID, e.g. 440601964878987871.
func (i ID) IsNaive() bool {
	l := len(i)
	if l < 18 {
		return false
	}
	var j int
	for j < l && '0' <= i[j] && i[j] <= '9' {
		j++
	}
	return j == l
}

const defaultSeparator = ":"

// IsComposited returns true if the given ID is a composited ID, e.g. x:y:z.
func (i ID) IsComposited(keyLength int) bool {
	if i.IsNaive() {
		return false
	}
	return i.matchKeyLength(keyLength)
}

func (i ID) matchKeyLength(l int) bool {
	switch l {
	case 0:
		return false
	case 1:
		return true
	}
	return strings.Count(string(i), defaultSeparator)+1 == l
}

// Split splits the composited key with default separator into an array,
// it should be called after IsNaive == false or IsComposited == true.
func (i ID) Split() []string {
	return strings.Split(string(i), defaultSeparator)
}
