package perspective

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/costreport"
	"github.com/seal-io/walrus/pkg/dao/types"
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
		values   []string
		pvalues  = make([]Value, 0)
		fieldStr = string(req.FieldName)
	)

	if req.FieldName.IsLabel() {
		/* Get the label values of given label name by the following sql:
		SELECT
		  DISTINCT(
		    labels ->> 'walrus.seal.io/project-name'
		  ) AS value
		FROM
		  "cost_reports"
		WHERE
		  "start_time" >= $1
		  AND "end_time" <= $2
		  AND "labels" <> 'null' :: jsonb
		*/
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

	if req.FieldName == types.FilterFieldConnectorID {
		/* Get the connector values of given connectorID field by the following sql:
		SELECT
		  name AS label,
		  id AS value
		FROM
		  "connectors"
		*/
		err := h.modelClient.Connectors().Query().
			Modify(func(s *sql.Selector) {
				s.
					Where(
						sql.EQ(connector.FieldCategory, types.ConnectorCategoryKubernetes)).
					SelectExpr(
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

	/* Get values of given field name by the following sql:
	SELECT
	  DISTINCT COALESCE(namespace, '')
	FROM
	  "cost_reports"
	WHERE
	  "start_time" >= '2023-09-09T00:00:00Z'
	  AND "end_time" <= '2023-09-15T16:16:51Z'
	*/

	fieldStr = fmt.Sprintf(`COALESCE(%s, '')`, req.FieldName)

	err := h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			if len(ps) != 0 {
				s.Where(sql.And(ps...))
			}
			s.Distinct().Select(fieldStr)
		}).
		Scan(req.Context, &values)
	if err != nil {
		return nil, 0, err
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
