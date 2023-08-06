package perspective

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/costreport"
	"github.com/seal-io/seal/pkg/dao/types"
)

func (h Handler) CollectionRouteGetFields(
	req CollectionRouteGetFieldsRequest,
) (CollectionRouteGetFieldsResponse, int, error) {
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
	case FieldTypeGroupBy:
		return builtinGroupFields, len(builtinGroupFields), nil
	case FieldTypeStep:
		return builtinStepFields, len(builtinStepFields), nil
	default:
		labelKeys, err := h.modelClient.CostReports().Query().
			Modify(func(s *sql.Selector) {
				s.Where(
					sql.And(ps...),
				).SelectExpr(
					sql.Expr(fmt.Sprintf(`DISTINCT(jsonb_object_keys(%s))`, costreport.FieldLabels)),
				)
			}).
			Strings(req.Context)
		if err != nil {
			return nil, 0, err
		}

		fields := builtinFilterFields
		for _, v := range labelKeys {
			fields = append(fields,
				Field{
					Label:     v,
					FieldName: types.LabelPrefix + v,
				},
			)
		}
		count := len(fields)

		return fields, count, nil
	}
}

func (h Handler) CollectionRouteGetFieldValues(
	req CollectionRouteGetFieldValuesRequest,
) (CollectionRouteGetFieldValuesResponse, int, error) {
	var ps []*sql.Predicate

	if req.StartTime != nil {
		ps = append(ps, sql.GTE(costreport.FieldStartTime, req.StartTime))
	}

	if req.EndTime != nil {
		ps = append(ps, sql.LTE(costreport.FieldEndTime, req.EndTime))
	}

	var (
		pvalues  = make([]Value, 0)
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
			}).
			Scan(req.Context, &s)
		if err != nil {
			return nil, 0, err
		}

		for _, v := range s {
			if v.Value == "" {
				continue
			}

			pvalues = append(pvalues, Value{
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
		}).
		Strings(req.Context)
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
			}).
			Scan(req.Context, &pvalues)
		if err != nil {
			return nil, 0, err
		}

		return pvalues, len(pvalues), nil
	}

	for _, v := range values {
		if v == "" {
			continue
		}

		pvalues = append(pvalues, Value{
			Label: v,
			Value: v,
		})
	}

	return pvalues, len(pvalues), nil
}
