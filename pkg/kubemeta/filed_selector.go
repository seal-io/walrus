package kubemeta

import (
	"strings"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/selection"
)

func FieldSelectorFromRequirements(requirements fields.Requirements) fields.Selector {
	s := make([]string, 0, len(requirements))
	for _, r := range requirements {
		switch r.Operator {
		case selection.DoubleEquals, selection.Equals, selection.NotEquals:
			s = append(s, r.Field+string(r.Operator)+fields.EscapeValue(r.Value))
		case selection.In, selection.NotIn:
			s = append(s, r.Field+string(r.Operator)+"("+r.Value+")")
		case selection.DoesNotExist:
			s = append(s, "!"+r.Field)
		case selection.Exists:
			s = append(s, r.Field)
		case selection.GreaterThan:
			s = append(s, r.Field+">"+r.Value)
		case selection.LessThan:
			s = append(s, r.Field+"<"+r.Value)
		}
	}
	fs, _ := fields.ParseSelector(strings.Join(s, ","))
	return fs
}
