package dao

import (
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
)

// Borrow From https://github.com/Masterminds/semver/blob/2f39fdc11c33c38e8b8b15b1f04334ba84e751f2/version.go#L42.
const semverExpression = `^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`

var OrderSemverVersionFunc = func(s *sql.Selector) {
	s.OrderExprFunc(func(b *sql.Builder) {
		b.WriteString("CASE WHEN")
		b.Ident(templateversion.FieldVersion)
		b.WriteString(" ~ '" + semverExpression + "' THEN string_to_array(regexp_replace(")
		b.Ident(templateversion.FieldVersion)
		b.WriteString(", E'[^0-9\\.]+','', 'g'), '.', '')::int[] ELSE NULL END DESC")
	})
}
