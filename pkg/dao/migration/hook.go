package migration

import (
	"context"

	"ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
)

// NotifySettingUpdate is a hook to create Trigger and Function to notify event while setting changed.
func NotifySettingUpdate() schema.ApplyHook {
	return func(applier schema.Applier) schema.Applier {
		return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
			plan.Changes = append(plan.Changes, []*migrate.Change{
				{
					Cmd: `CREATE OR REPLACE FUNCTION process_setting_update() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'UPDATE') THEN
        PERFORM pg_notify(
                'setting_update',
                json_build_object(
                    'tableName', TG_TABLE_NAME,
                    'operation', TG_OP, 'recordID', COALESCE(NEW.id, 0)
                    )::text
            );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;`,
				},
				{
					Cmd: `CREATE OR REPLACE TRIGGER setting_update_trigger
    BEFORE UPDATE
    ON settings
    FOR EACH ROW
EXECUTE FUNCTION process_setting_update();`,
				},
			}...)

			return applier.Apply(ctx, conn, plan)
		})
	}
}
