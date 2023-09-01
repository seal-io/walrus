package log

// Filter holds the operations to filter out the given log arguments.
type Filter func(msg string, keysAndValues []any) (string, []any, bool)
