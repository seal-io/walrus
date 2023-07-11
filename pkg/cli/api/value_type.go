package api

// Value types extend the OpenAPI specification to support more types.
const (
	ValueTypeObjectID        = "objectID"
	ValueTypeArrayObject     = "array[object]"
	ValueTypeArrayString     = "array[string]"
	ValueTypeArrayNumber     = "array[number]"
	ValueTypeArrayInt        = "array[integer]"
	ValueTypeArrayBoolean    = "array[boolean]"
	ValueTypeMapStringInt64  = "map[string]int64"
	ValueTypeMapStringInt32  = "map[string]int32"
	ValueTypeMapStringInt    = "map[string]int"
	ValueTypeMapStringString = "map[string]string"
)
