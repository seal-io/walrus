package intercept

import "entgo.io/ent/dialect/sql"

type (
	instructionKey   string
	instructionValue struct{}
)

type target interface {
	WhereP(...func(*sql.Selector))
}
