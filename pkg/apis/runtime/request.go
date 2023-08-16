package runtime

import (
	"strings"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/strs"
)

// RequestCollection holds the requesting data of collection,
// including querying, sorting, extracting and pagination.
type RequestCollection[Q, S ~func(*sql.Selector)] struct {
	RequestQuerying[Q] `query:",inline"`
	RequestSorting[S]  `query:",inline"`
	RequestExtracting  `query:",inline"`
	RequestPagination  `query:",inline"`
}

// RequestQuerying holds the requesting query data.
type RequestQuerying[T ~func(s *sql.Selector)] struct {
	// Query specifies the content to search some preset fields,
	// it's a case-insensitive fuzzy filter,
	// i.e. /v1/repositories?query=repo%2Fname.
	Query *string `query:"query,omitempty"`
}

// Querying returns an OR predicate with the given search fields,
// returns false if there is no query requesting.
func (r RequestQuerying[T]) Querying(searchFields []string) (T, bool) {
	if r.Query == nil || len(searchFields) == 0 {
		return nil, false
	}

	p := func(s *sql.Selector) {
		q := make([]*sql.Predicate, 0, len(searchFields))
		for _, f := range searchFields {
			q = append(q, sql.ContainsFold(s.C(f), *r.Query))
		}

		if len(q) == 1 {
			s.Where(q[0])
		} else {
			s.Where(sql.Or(q...))
		}
	}

	return p, true
}

// RequestSorting holds the requesting sort data.
type RequestSorting[T ~func(*sql.Selector)] struct {
	// Sorts specifies the fields for sorting,
	// i.e. /v1/repositories?sort=-createTime&sort=name.
	Sorts []string `query:"sort,omitempty"`
}

// WithAsc appends the asc sorting field list to the sorting list.
func (r RequestSorting[T]) WithAsc(fields ...string) RequestSorting[T] {
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			continue
		}

		r.Sorts = append(r.Sorts, fields[i])
	}

	return r
}

// WithDesc appends the desc sorting list to the sorting list.
func (r RequestSorting[T]) WithDesc(fields ...string) RequestSorting[T] {
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			continue
		}

		r.Sorts = append(r.Sorts, "-"+fields[i])
	}

	return r
}

// Sorting returns the order list with the given allow list,
// returns false if there are not any sorting key requesting and default list.
func (r RequestSorting[T]) Sorting(allowKeys []string, defaultOrders ...T) ([]T, bool) {
	if len(r.Sorts) == 0 || len(allowKeys) == 0 {
		return defaultOrders, len(defaultOrders) != 0
	}

	orders := make([]T, 0, len(allowKeys))
	allows := sets.NewString(allowKeys...)

	for i := 0; i < len(r.Sorts); i++ {
		key := strings.TrimSpace(r.Sorts[i])
		if key == "" {
			continue
		}

		order := model.Asc

		switch key[0] {
		case '-':
			order = model.Desc
			key = key[1:]
		case '+':
			key = key[1:]
		}

		if allows.Has(key) {
			allows.Delete(key) // Not allow duplicate inputs.
			orders = append(orders, order(key))
		} else if ukey := strs.Underscore(key); ukey != key {
			if allows.Has(ukey) {
				allows.Delete(ukey) // Not allow duplicate inputs.
				orders = append(orders, order(ukey))
			}
		}
	}

	if len(orders) == 0 {
		return defaultOrders, len(defaultOrders) != 0
	}

	return orders, true
}

// RequestExtracting holds the requesting extraction data.
type RequestExtracting struct {
	// Extracts specifies the fields for querying,
	// i.e. /v1/repositories?extract=-id&extract=name.
	Extracts []string `query:"extract,omitempty"`
}

// With appends the included field list to the extracting list.
func (r RequestExtracting) With(fields ...string) RequestExtracting {
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			continue
		}

		r.Extracts = append(r.Extracts, fields[i])
	}

	return r
}

// Without appends the excluded field list to the extracting list.
func (r RequestExtracting) Without(fields ...string) RequestExtracting {
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			continue
		}

		r.Extracts = append(r.Extracts, "-"+fields[i])
	}

	return r
}

// Extracting returns the field list with the given allow list,
// returns false if there are not any extracting key requesting and default list.
func (r RequestExtracting) Extracting(allowFields []string, defaultFields ...string) ([]string, bool) {
	if len(r.Extracts) == 0 || len(allowFields) == 0 {
		return defaultFields, len(defaultFields) != 0
	}

	candidates := make([]string, len(r.Extracts)+len(defaultFields))
	copy(candidates, r.Extracts)
	copy(candidates[len(r.Extracts):], defaultFields)

	fields := make([]string, 0, len(candidates))
	allows := sets.NewString(allowFields...)

	for i := 0; i < len(candidates); i++ {
		key := strings.TrimSpace(candidates[i])
		if key == "" {
			continue
		}

		with := true

		switch candidates[i][0] {
		case '-':
			with = false
			key = candidates[i][1:]
		case '+':
			key = candidates[i][1:]
		}

		if allows.Has(key) {
			allows.Delete(key) // Not allow duplicate inputs.

			if with {
				fields = append(fields, key)
			}
		} else if ukey := strs.Underscore(key); ukey != key {
			if allows.Has(ukey) {
				allows.Delete(ukey) // Not allow duplicate inputs.

				if with {
					fields = append(fields, ukey)
				}
			}
		}
	}

	if len(fields) == 0 {
		return defaultFields, len(defaultFields) != 0
	}

	return fields, true
}

// RequestPagination holds the requesting pagination data.
type RequestPagination struct {
	// Page specifies the page number for querying,
	// i.e. /v1/repositories?page=1&perPage=10.
	Page int `query:"page,default=1"`

	// PerPage specifies the page size for querying,
	// i.e. /v1/repositories?page=1&perPage=10.
	PerPage int `query:"perPage,default=100"`
}

// Limit returns the limit of paging.
func (r RequestPagination) Limit() int {
	limit := r.PerPage
	if limit <= 0 {
		limit = 100
	}

	return limit
}

// Offset returns the offset of paging.
func (r RequestPagination) Offset() int {
	offset := r.Limit() * (r.Page - 1)
	if offset < 0 {
		offset = 0
	}

	return offset
}

// Paging returns the limit and offset of paging,
// returns false if there is no pagination requesting.
func (r RequestPagination) Paging() (limit, offset int, request bool) {
	request = r.Page > 0
	if !request {
		return
	}
	limit = r.Limit()
	offset = r.Offset()

	return
}
