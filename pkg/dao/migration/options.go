package migration

import "entgo.io/ent/dialect/sql/schema"

// Options indicate the migration options.
var Options = []schema.MigrateOption{
	schema.WithApplyHook(NotifySettingUpdate()),
}
