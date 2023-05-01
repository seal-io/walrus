package runtime

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gorilla/websocket"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/strs"
)

type RequestCollection[T ~func(*sql.Selector)] struct {
	RequestQuerying[T] `query:",inline"`
	RequestSorting     `query:",inline"`
	RequestPagination  `query:",inline"`
	RequestExtracting  `query:",inline"`
}

func (r RequestCollection[T]) Validate() (err error) {
	var validates = []func() error{
		r.RequestQuerying.Validate,
		r.RequestSorting.Validate,
		r.RequestExtracting.Validate,
	}
	for i := range validates {
		if multierr.AppendInto(&err, validates[i]()) {
			return
		}
	}
	return
}

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
	var limit = r.PerPage
	if limit <= 0 {
		limit = 100
	}
	return limit
}

// Offset returns the offset of paging.
func (r RequestPagination) Offset() int {
	var offset = r.Limit() * (r.Page - 1)
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

type RequestSorting struct {
	// Sorts specifies the fields for sorting,
	// i.e. /v1/repositories?sort=-createTime&sort=name.
	Sorts []string `query:"sort,omitempty"`
}

func (r RequestSorting) Validate() error {
	for i := 0; i < len(r.Sorts); i++ {
		if strings.TrimSpace(r.Sorts[i]) == "" {
			return errors.New("blank sort value is not allowed")
		}
	}
	return nil
}

// WithAsc appends the asc sorting field list to the sorting list.
func (r RequestSorting) WithAsc(fields ...string) RequestSorting {
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			continue
		}
		r.Sorts = append(r.Sorts, fields[i])
	}
	return r
}

// WithDesc appends the desc sorting list to the sorting list.
func (r RequestSorting) WithDesc(fields ...string) RequestSorting {
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
func (r RequestSorting) Sorting(allowKeys []string, defaultOrders ...model.OrderFunc) ([]model.OrderFunc, bool) {
	if len(r.Sorts) == 0 || len(allowKeys) == 0 {
		return defaultOrders, len(defaultOrders) != 0
	}

	var orders = make([]model.OrderFunc, 0, len(allowKeys))
	var allows = sets.NewString(allowKeys...)
	for i := 0; i < len(r.Sorts); i++ {
		if r.Sorts[i] == "" {
			continue
		}
		var order = model.Asc
		var key string
		switch r.Sorts[i][0] {
		case '-':
			order = model.Desc
			key = r.Sorts[i][1:]
		case '+':
			key = r.Sorts[i][1:]
		default:
			key = r.Sorts[i]
		}
		if allows.Has(key) {
			allows.Delete(key) // not allow duplicate inputs
			orders = append(orders, order(key))
		} else if ukey := strs.Underscore(key); ukey != key {
			if allows.Has(ukey) {
				allows.Delete(ukey) // not allow duplicate inputs
				orders = append(orders, order(ukey))
			}
		}
	}
	if len(orders) == 0 {
		return defaultOrders, len(defaultOrders) != 0
	}
	return orders, true
}

type RequestGrouping struct {
	// Groups specifies the fields for grouping,
	// i.e. /v1/repositories?group=namespace&group=name.
	Groups []string `query:"group,omitempty"`
}

func (r RequestGrouping) Validate() error {
	if len(r.Groups) > 3 {
		return errors.New("too deep in group levels")
	}
	for i := 0; i < len(r.Groups); i++ {
		if strings.TrimSpace(r.Groups[i]) == "" {
			return errors.New("blank group value is not allowed")
		}
	}
	return nil
}

type RequestExtracting struct {
	// Extracts specifies the fields for querying,
	// i.e. /v1/repositories?extract=-id&extract=name.
	Extracts []string `query:"extract,omitempty"`
}

func (r RequestExtracting) Validate() error {
	for i := 0; i < len(r.Extracts); i++ {
		if strings.TrimSpace(r.Extracts[i]) == "" {
			return errors.New("blank extract value is not allowed")
		}
	}
	return nil
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

	var candidates = make([]string, len(r.Extracts)+len(defaultFields))
	copy(candidates, r.Extracts)
	copy(candidates[len(r.Extracts):], defaultFields)

	var fields = make([]string, 0, len(candidates))
	var allows = sets.NewString(allowFields...)
	for i := 0; i < len(candidates); i++ {
		if candidates[i] == "" {
			continue
		}
		var with = true
		var key string
		switch candidates[i][0] {
		case '-':
			with = false
			key = candidates[i][1:]
		case '+':
			key = candidates[i][1:]
		default:
			key = candidates[i]
		}
		if allows.Has(key) {
			allows.Delete(key) // not allow duplicate inputs
			if with {
				fields = append(fields, key)
			}
		} else if ukey := strs.Underscore(key); ukey != key {
			if allows.Has(ukey) {
				allows.Delete(ukey) // not allow duplicate inputs
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

// ExtractingSet is similar to Extracting but returns a sets.Set[string] of fields.
func (r RequestExtracting) ExtractingSet(allowFields []string, defaultFields ...string) sets.Set[string] {
	var fields, ok = r.Extracting(allowFields, defaultFields...)
	if !ok {
		return sets.Set[string]{}
	}
	return sets.New[string](fields...)
}

type RequestQuerying[T ~func(s *sql.Selector)] struct {
	// Query specifies the content to search some preset fields,
	// it's a case-insenstive fuzzy filter,
	// i.e. /v1/repositories?query=repo%2Fname.
	Query *string `query:"query,omitempty"`
}

func (r RequestQuerying[T]) Validate() error {
	if r.Query != nil && strings.TrimSpace(*r.Query) == "" {
		return errors.New("blank query value is not allowed")
	}
	return nil
}

// Querying returns an OR predicate with the given search fields,
// returns false if there is no query requesting.
func (r RequestQuerying[T]) Querying(searchFields []string) (T, bool) {
	if r.Query == nil || len(searchFields) == 0 {
		return nil, false
	}
	var p = func(s *sql.Selector) {
		var q = make([]*sql.Predicate, 0, len(searchFields))
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

type RequestOperating struct {
	// Action specifies the action for operating,
	// i.e. /v1/users/:id/logs?action=count.
	Action *string `query:"action,omitempty"`
}

func (r RequestOperating) Validate() error {
	if r.Action != nil && strings.TrimSpace(*r.Action) == "" {
		return errors.New("blank action value is not allowed")
	}
	return nil
}

type RequestStream struct {
	firstReadOnce *sync.Once
	firstReadChan <-chan struct {
		t int
		r io.Reader
		e error
	}
	ctx       context.Context
	ctxCancel func()
	conn      *websocket.Conn
}

// SendMsg sends the given data to client.
func (r RequestStream) SendMsg(data []byte) error {
	var _, err = r.Write(data)
	return err
}

// SendJSON marshals the given object as JSON and sends to client.
func (r RequestStream) SendJSON(i any) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return r.SendMsg(bs)
}

// RecvMsg receives message from client.
func (r RequestStream) RecvMsg() ([]byte, error) {
	return io.ReadAll(r)
}

// RecvJSON receives JSON message from client and unmarshals into the given object.
func (r RequestStream) RecvJSON(i any) error {
	bs, err := r.RecvMsg()
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, i)
}

// Write implements io.Writer.
func (r RequestStream) Write(p []byte) (n int, err error) {
	msgWriter, err := r.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	defer func() { _ = msgWriter.Close() }()
	return msgWriter.Write(p)
}

// Read implements io.Reader.
func (r RequestStream) Read(p []byte) (n int, err error) {
	var (
		firstRead bool
		msgType   int
		msgReader io.Reader
	)
	r.firstReadOnce.Do(func() {
		firstRead = true
		var fr = <-r.firstReadChan
		msgType, msgReader, err = fr.t, fr.r, fr.e
	})
	if !firstRead {
		msgType, msgReader, err = r.conn.NextReader()
	}
	if err != nil {
		return
	}
	switch msgType {
	default:
		err = &websocket.CloseError{
			Code: websocket.CloseUnsupportedData,
			Text: "unresolved message type: binary",
		}
		return
	case websocket.TextMessage:
	}
	return msgReader.Read(p)
}

// Cancel cancels the underlay context.Context.
func (r RequestStream) Cancel() {
	r.ctxCancel()
}

// Deadline implements context.Context.
func (r RequestStream) Deadline() (deadline time.Time, ok bool) {
	return r.ctx.Deadline()
}

// Done implements context.Context.
func (r RequestStream) Done() <-chan struct{} {
	return r.ctx.Done()
}

// Err implements context.Context.
func (r RequestStream) Err() error {
	return r.ctx.Err()
}

// Value implements context.Context.
func (r RequestStream) Value(key any) any {
	return r.ctx.Value(key)
}
