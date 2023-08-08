package perspective

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/perspective/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/costreport"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/types"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "Perspective"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	entity, err := h.modelClient.Perspectives().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposePerspective(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Perspectives().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.Perspectives().UpdateOne(entity).
		Set(entity).
		Exec(ctx)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	entity, err := h.modelClient.Perspectives().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposePerspective(entity), nil
}

// Batch APIs.

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Perspectives().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

var (
	queryFields = []string{
		perspective.FieldName,
	}
	getFields  = perspective.Columns
	sortFields = []string{
		perspective.FieldName,
		perspective.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Perspectives().Query().
		Where(perspective.NameContains(req.Name))
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(perspective.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposePerspectives(entities), cnt, nil
}

// Extensional APIs.

func (h Handler) CollectionRouteFields(
	ctx *gin.Context,
	req view.CollectionRouteFieldsRequest,
) (view.CollectionRouteFieldsResponse, int, error) {
	ps := []*sql.Predicate{
		sqljson.ValueIsNotNull(costreport.FieldLabels),
	}
	if req.StartTime != nil {
		ps = append(ps, sql.GTE(costreport.FieldStartTime, req.StartTime))
	}

	if req.EndTime != nil {
		ps = append(ps, sql.LTE(costreport.FieldEndTime, req.EndTime))
	}

	switch req.FieldType {
	case view.FieldTypeGroupBy:
		return view.BuiltInPerspectiveGroupFields, len(view.BuiltInPerspectiveGroupFields), nil
	case view.FieldTypeStep:
		return view.BuiltInPerspectiveStepFields, len(view.BuiltInPerspectiveStepFields), nil
	default:
		labelKeys, err := h.modelClient.CostReports().Query().
			Modify(func(s *sql.Selector) {
				s.Where(
					sql.And(ps...),
				).SelectExpr(
					sql.Expr(fmt.Sprintf(`DISTINCT(jsonb_object_keys(%s))`, costreport.FieldLabels)),
				)
			}).
			Strings(ctx)
		if err != nil {
			return nil, 0, err
		}

		fields := view.BuiltInPerspectiveFilterFields
		for _, v := range labelKeys {
			fields = append(fields, view.LabelKeyToPerspectiveField(v))
		}
		count := len(fields)

		return fields, count, nil
	}
}

func (h Handler) CollectionRouteValues(
	ctx *gin.Context,
	req view.CollectionRouteFieldValuesRequest,
) (view.CollectionRouteValuesResponse, int, error) {
	var ps []*sql.Predicate
	if req.StartTime != nil {
		ps = append(ps, sql.GTE(costreport.FieldStartTime, req.StartTime))
	}

	if req.EndTime != nil {
		ps = append(ps, sql.LTE(costreport.FieldEndTime, req.EndTime))
	}

	var (
		pvalues  = make([]view.PerspectiveValue, 0)
		fieldStr = string(req.FieldName)
	)

	if req.FieldName.IsLabel() {
		var (
			s []struct {
				Value string `json:"value"`
			}
			field = fmt.Sprintf(`labels ->> '%s'`,
				strings.TrimPrefix(fieldStr, types.LabelPrefix))
		)

		ps = append(ps, sqljson.ValueIsNotNull(costreport.FieldLabels))

		err := h.modelClient.CostReports().Query().
			Modify(func(s *sql.Selector) {
				s.Where(
					sql.And(ps...),
				).SelectExpr(
					sql.Expr(fmt.Sprintf(`DISTINCT(%s) AS value`, field)),
				)
			}).Scan(ctx, &s)
		if err != nil {
			return nil, 0, err
		}

		for _, v := range s {
			if v.Value == "" {
				continue
			}

			pvalues = append(pvalues, view.PerspectiveValue{
				Label: v.Value,
				Value: v.Value,
			})
		}

		return pvalues, len(pvalues), nil
	}

	values, err := h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			if len(ps) != 0 {
				s.Where(sql.And(ps...))
			}
			s.Distinct().Select(fieldStr)
		}).Strings(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.FieldName == types.FilterFieldConnectorID {
		var ids []any
		for _, v := range values {
			ids = append(ids, v)
		}

		err = h.modelClient.Connectors().Query().
			Modify(func(s *sql.Selector) {
				s.Where(
					sql.In(connector.FieldID, ids...),
				).SelectExpr(
					sql.Expr(fmt.Sprintf(`%s AS label`, connector.FieldName)),
					sql.Expr(fmt.Sprintf(`%s AS value`, connector.FieldID)),
				)
			}).Scan(ctx, &pvalues)
		if err != nil {
			return nil, 0, err
		}

		return pvalues, len(pvalues), nil
	}

	for _, v := range values {
		if v == "" {
			continue
		}

		pvalues = append(pvalues, view.PerspectiveValue{
			Label: v,
			Value: v,
		})
	}

	return pvalues, len(pvalues), nil
}
