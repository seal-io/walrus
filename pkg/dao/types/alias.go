package types

import (
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ID shows the primary key in string but stores in big integer,
// also be good at catching composited primary keys.
type ID = oid.ID
