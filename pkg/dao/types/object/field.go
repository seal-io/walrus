package object

import (
	"database/sql/driver"
	"errors"
	"strconv"
)

// ID shows the primary key in string but stores in big integer,
// also be good at catching composited primary keys.
type ID string

// NewID creates an ID with the given integer.
func NewID(id uint64) ID {
	return ID(strconv.FormatUint(id, 10))
}

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

// Valid returns true if the given ID is a naive numeric value at least 18 digits,
// e.g. 440601964878987871.
func (i ID) Valid() bool {
	return len(i) >= 18 && isNumeric(string(i))
}

func isNumeric(s string) bool {
	var j int
	for j < len(s) && '0' <= s[j] && s[j] <= '9' {
		j++
	}

	return j == len(s)
}
