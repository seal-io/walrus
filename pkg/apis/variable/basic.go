package variable

import (
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/variable"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	entity, err := h.modelClient.Variables().Create().
		Set(entity).
		Save(req.Context)
	if err != nil {
		return nil, err
	}

	return exposeVariable(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.Variables().UpdateOne(entity).
		Set(entity).
		Exec(req.Context)
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.Variables().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		variable.FieldName,
	}
	getFields = variable.WithoutFields(
		variable.FieldUpdateTime)
	sortFields = []string{
		variable.FieldName,
		variable.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	/* SQL generate for count.

	SELECT
	  COUNT("variables"."id")
	FROM
	  "variables"
	  JOIN (
	    SELECT
	      "alias"."name",
	      MAX(
	        CASE WHEN project_id IS NOT NULL
	        AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	        AND environment_id IS NULL THEN 2 ELSE 1 END
	      ) AS scope
	    FROM
	      "variables" AS "alias"
	    WHERE
	      (
	        "project_id" = $1
	        AND "environment_id" = $2
	      )
	      OR (
	        "project_id" = $3
	        AND "environment_id" IS NULL
	      )
	      OR "project_id" IS NULL
	    GROUP BY
	      "alias"."name"
	  ) AS "alias" ON (
	    (
	      "project_id" = $4
	      AND "environment_id" = $5
	    )
	    OR (
	      "project_id" = $6
	      AND "environment_id" IS NULL
	    )
	    OR "project_id" IS NULL
	  )
	  AND "variables"."name" = "alias"."name"
	  AND (
	    CASE WHEN project_id IS NOT NULL
	    AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	    AND environment_id IS NULL THEN 2 ELSE 1 END
	  ) = alias.scope
	*/
	/* SQL generate for query.

	SELECT
	  "variables"."id",
	  "variables"."create_time",
	  "variables"."project_id",
	  "variables"."name",
	  "variables"."value",
	  "variables"."sensitive",
	  "variables"."description",
	  "variables"."environment_id"
	FROM
	  "variables"
	  JOIN (
	    SELECT
	      "alias"."name",
	      MAX(
	        CASE WHEN project_id IS NOT NULL
	        AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	        AND environment_id IS NULL THEN 2 ELSE 1 END
	      ) AS scope
	    FROM
	      "variables" AS "alias"
	    WHERE
	      (
	        (
	          "project_id" = $1
	          AND "environment_id" = $2
	        )
	        OR (
	          "project_id" = $3
	          AND "environment_id" IS NULL
	        )
	        OR "project_id" IS NULL
	      )
	    GROUP BY
	      "alias"."name"
	  ) AS "alias" ON (
	    (
	      "project_id" = $4
	      AND "environment_id" = $5
	    )
	    OR (
	      "project_id" = $6
	      AND "environment_id" IS NULL
	    )
	    OR "project_id" IS NULL
	  )
	  AND "variables"."name" = "alias"."name"
	  AND (
	    CASE WHEN project_id IS NOT NULL
	    AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	    AND environment_id IS NULL THEN 2 ELSE 1 END
	  ) = alias.scope
	ORDER BY
	  "variables"."create_time" DESC
	LIMIT
	  100
	*/
	// Ps is the generated query condition base on input.
	var ps *sql.Predicate

	switch {
	case req.Environment != nil:
		// Environment scope only.
		ps = sql.And(
			sql.EQ(variable.FieldProjectID, req.Environment.Project.ID),
			sql.EQ(variable.FieldEnvironmentID, req.Environment.ID),
		)

		if req.IncludeInherited {
			ps = sql.Or(
				// Global scope.
				sql.IsNull(variable.FieldProjectID),
				// Project scope.
				sql.And(
					sql.EQ(variable.FieldProjectID, req.Environment.Project.ID),
					sql.IsNull(variable.FieldEnvironmentID),
				),
				// Environment scope.
				ps,
			)
		}
	case req.Project != nil:
		// Project scope only.
		ps = sql.And(
			sql.EQ(variable.FieldProjectID, req.Project.ID),
			sql.IsNull(variable.FieldEnvironmentID),
		)

		if req.IncludeInherited {
			ps = sql.Or(
				// Global scope.
				sql.IsNull(variable.FieldProjectID),
				// Project scope.
				ps,
			)
		}
	default:
		// Global scope.
		ps = sql.IsNull(variable.FieldProjectID)
	}

	scopeExpr := `CASE WHEN project_id IS NOT NULL AND environment_id IS NOT NULL THEN 3
                                          WHEN project_id IS NOT NULL AND environment_id IS NULL THEN 2
                                          ELSE 1 END`

	modifier := func(s *sql.Selector) {
		alias := sql.Table(variable.Table).
			As("alias")

		subQuery := sql.Select(alias.C(variable.FieldName)).
			AppendSelectExpr(sql.Expr(fmt.Sprintf("MAX(%s) AS scope", scopeExpr))).
			From(alias).
			Where(ps).
			GroupBy(alias.C(variable.FieldName)).
			As("alias")

		s.Join(subQuery).
			OnP(
				sql.And(
					ps,
					sql.ColumnsEQ(
						s.C(variable.FieldName),
						alias.C(variable.FieldName),
					),
					sql.ExprP(fmt.Sprintf("(%s) = alias.scope", scopeExpr)),
				),
			)
	}

	// Generate query.
	query := h.modelClient.Variables().Query().
		Modify(modifier)

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().
		Modify(modifier).
		Count(req.Context)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(variable.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return exposeVariables(entities), cnt, err
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.Variables().Delete().
			Where(variable.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
