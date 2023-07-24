package dashboard

import (
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/types/object"
)

// predicateIn converts the given field name and values into an IN predicate.
func predicateIn[T ~func(*sql.Selector)](isAdmin bool, name string, values []object.ID) []T {
	if isAdmin {
		return nil
	}

	// NB(thxCode): For non-admin,
	// an empty value list will get a FALSE clause.
	return []T{
		sql.FieldIn(name, values...),
	}
}

// predicateOr is different from the commonly used OR predicate,
// it returns nil if given others argument is empty.
func predicateOr[T ~func(*sql.Selector)](first T, others ...T) []T {
	if len(others) == 0 {
		return nil
	}

	return []T{
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			first(s1)
			for _, p := range others {
				s1.Or()
				p(s1)
			}
			s.Where(s1.P())
		},
	}
}
