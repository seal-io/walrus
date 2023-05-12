package session

import "net/http"

type operator = uint8

const (
	operatingGet    operator = 0x01
	operatingPost   operator = 0x02
	operatingPut    operator = 0x04
	operatingDelete operator = 0x08
	operatingAny             = operatingGet | operatingPost | operatingPut | operatingDelete
)

func getOperators() map[operator]int {
	return map[operator]int{
		operatingGet:    0,
		operatingPost:   1,
		operatingPut:    2,
		operatingDelete: 3,
	}
}

// Permission holds the all operations of permission.
type Permission [4]Operation

// If returns the Operation with the given operator.
func (p Permission) If(s string) (pv Operation) {
	var pk operator
	switch s {
	case http.MethodPost:
		pk = operatingPost
	case http.MethodDelete:
		pk = operatingDelete
	case http.MethodPut:
		pk = operatingPut
	case http.MethodGet:
		pk = operatingGet
	}
	if pk == 0 {
		return
	}
	var idx = getOperators()[pk]
	return p[idx]
}

// Operation holds the operation of a specified permission,
// the usage is as below:
//   - {scope: "private",includes: nil, excludes: nil}: operate with private thing.
//   - {scope: "shared", includes: nil, excludes: nil}: operate with shared thing.
//   - {scope: "",       includes: [.], excludes: nil}: operate with something included.
//   - {scope: "",       includes: nil, excludes: [.]}: operate with something not excluded.
//   - {scope: "",       includes: nil, excludes: nil}: operate with anything.
type Operation struct {
	work     bool
	scope    string
	includes []string
	excludes []string
}

type OperateTarget = uint8

const (
	None OperateTarget = iota + 1
	Private
	Shared
	Something
	Any
)

// Then return the OperateTarget of this Operation.
func (p Operation) Then() (OperateTarget, []string, []string) {
	if !p.work {
		return None, nil, nil
	}
	switch p.scope {
	case "private":
		return Private, nil, nil
	case "shared":
		return Shared, nil, nil
	default:
	}
	if len(p.includes) != 0 || len(p.excludes) != 0 {
		return Something, p.includes, p.excludes
	}
	return Any, nil, nil
}

func (p Operation) merge(v Operation) (o Operation) {
	if !p.work {
		v.work = true
		return v
	}

	var pf, pi, pe = p.Then()
	var vf, vi, ve = v.Then()

	switch {
	case vf < pf:
		return p
	case vf > pf:
		return v
	case vf != Something:
		return v
	}

	if len(pi) == 0 {
		// Excludes in p.
		if len(vi) == 0 {
			// Excludes in v, merge both.
			o.excludes = make([]string, 0, len(pe)+len(ve))
			o.excludes = append(o.excludes, pe...)
			o.excludes = append(o.excludes, ve...)
		} else {
			// Includes in v, excludes get higher priority.
			o.excludes = pe
		}
	} else {
		// Includes in p.
		if len(vi) == 0 {
			// Excludes in v, excludes get higher priority.
			o.excludes = ve
		} else {
			// Includes in v, merge both.
			o.includes = make([]string, 0, len(pi)+len(vi))
			o.includes = append(o.includes, pi...)
			o.includes = append(o.includes, vi...)
		}
	}
	return
}
