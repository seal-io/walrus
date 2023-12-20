package migration

import (
	"context"
	"text/template"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
)

const (
	modelChangeFunction = "model_change"
	ModelChangeChannel  = modelChangeFunction + "_channel"
)

// ApplyModelChangeTrigger returns the schema.MigrateOption slice for
// creating the model change trigger,
// which is used to notify the server with ModelChangeChannel when the database is changed.
func ApplyModelChangeTrigger(tableNames []string) (opts []schema.MigrateOption) {
	hooks := []schema.ApplyHook{
		createModelChangeFunction(),
		createModelChangeTrigger(tableNames),
	}

	for i := range hooks {
		opts = append(opts, schema.WithApplyHook(hooks[i]))
	}

	return
}

func createModelChangeFunction() schema.ApplyHook {
	const tmpl = `
CREATE OR REPLACE FUNCTION {{ .Function }}() RETURNS TRIGGER
    LANGUAGE plpgsql
AS 
$$
DECLARE
    ids                 jsonb;
    ids_length          int;
    column_sum        int;
    payload             text;
BEGIN

    -- Check if the table has project_id and environment_id columns. 
    -- For entity tables, project_id is always present when environment_id is present.
    SELECT SUM(
      CASE
        WHEN column_name = 'name' then 100
        WHEN column_name = 'project_id' then 10
        WHEN column_name = 'environment_id' then 1
        ELSE 0
      END)
    FROM information_schema.columns
    WHERE table_name = TG_TABLE_NAME
    INTO column_sum;

    -- Build the ID list.
    ids := '[]'::jsonb;

    CASE TG_OP
    WHEN 'INSERT', 'UPDATE' THEN
        EXECUTE 'SELECT jsonb_agg(jsonb_build_object(''id'', id::text' ||
            CASE WHEN column_sum >= 100 THEN ', ''name'', name::text' ELSE '' END ||
            CASE WHEN column_sum % 100 >= 10 THEN ', ''project_id'', project_id::text' ELSE '' END ||
            CASE WHEN column_sum % 10 >= 1 THEN ', ''environment_id'', environment_id::text' ELSE '' END ||
            ')) FROM new_table' INTO ids;
    WHEN 'DELETE' THEN
        EXECUTE 'SELECT jsonb_agg(jsonb_build_object(''id'', id::text' ||
            CASE WHEN column_sum >= 100 THEN ', ''name'', name::text' ELSE '' END ||
            CASE WHEN column_sum % 100 >= 10 THEN ', ''project_id'', project_id::text' ELSE '' END ||
            CASE WHEN column_sum % 10 >= 1 THEN ', ''environment_id'', environment_id::text' ELSE '' END ||
            ')) FROM old_table' INTO ids;
    ELSE
        RAISE EXCEPTION 'Unknown Operation';
    END CASE;

    -- Validate the length of ID list.
    ids_length := jsonb_array_length(ids);
    IF (ids_length = 0) THEN
        RETURN NULL;
    END IF;

    -- Build the notification payload.
    payload := json_build_object(
        'ts', current_timestamp,
        'op', lower(TG_OP),
        'tb_s', TG_TABLE_SCHEMA,
        'tb_n', TG_TABLE_NAME,
        'ids', ids)::text;

    -- Notify the channel.
    PERFORM pg_notify('{{ .Channel }}', payload);

    RETURN NULL;
END;
$$;
`

	tpl := template.Must(template.New("tmpl").Parse(tmpl))

	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			plan.Changes = append(plan.Changes, &migrate.Change{
				Cmd: executeTemplate(tpl, map[string]any{
					"Function": modelChangeFunction,
					"Channel":  ModelChangeChannel,
				}),
			})

			return next.Apply(ctx, conn, plan)
		})
	}
}

func createModelChangeTrigger(tableNames []string) schema.ApplyHook {
	const tmpl = `
CREATE OR REPLACE TRIGGER {{ .Function }}_ins
	AFTER INSERT ON {{ .Table }}
	REFERENCING NEW TABLE AS new_table
	FOR EACH STATEMENT EXECUTE FUNCTION {{ .Function }}();
CREATE OR REPLACE TRIGGER {{ .Function }}_upd
	AFTER UPDATE ON {{ .Table }}
	REFERENCING OLD TABLE AS old_table NEW TABLE AS new_table
	FOR EACH STATEMENT EXECUTE FUNCTION {{ .Function }}();
CREATE OR REPLACE TRIGGER {{ .Function }}_del
	AFTER DELETE ON {{ .Table }}
	REFERENCING OLD TABLE AS old_table
	FOR EACH STATEMENT EXECUTE FUNCTION {{ .Function }}();
`

	tpl := template.Must(template.New("tmpl").Parse(tmpl))

	return func(next schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			for i := range tableNames {
				plan.Changes = append(plan.Changes, &migrate.Change{
					Cmd: executeTemplate(tpl, map[string]any{
						"Function": modelChangeFunction,
						"Table":    tableNames[i],
					}),
				})
			}

			return next.Apply(ctx, conn, plan)
		})
	}
}
