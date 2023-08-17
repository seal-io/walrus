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
    ids        json;
    ids_length integer;
    payload    text;
BEGIN
    
    -- Select ID list from transition table,
    -- and convert the result to a JSON string array.
    CASE TG_OP
    WHEN 'INSERT' THEN
        SELECT array_to_json(array(SELECT id::text FROM new_table)) INTO ids;
    WHEN 'UPDATE' THEN
        SELECT array_to_json(array(SELECT id::text FROM new_table)) INTO ids;
    WHEN 'DELETE' THEN
        SELECT array_to_json(array(SELECT id::text FROM old_table)) INTO ids;
    ELSE
        RAISE EXCEPTION 'Unknown Operation';
    END CASE;
    
    -- Validate the length of ID list.
    ids_length := json_array_length(ids);
    IF (ids_length = 0) THEN
        RETURN NULL;
    END IF;

    -- Build the notification payload.
    payload := json_build_object(
        'ts',current_timestamp,
        'op',lower(TG_OP),
        'tb_s',TG_TABLE_SCHEMA,
        'tb_n',TG_TABLE_NAME,
        'ids',ids)::text;
    
    -- Notify the channel.
    PERFORM pg_notify('{{ .Channel }}',payload);

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
