package json

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Get returns gjson.Result after searching the specified path in given JSON data,
// the path rule is documented on https://github.com/tidwall/gjson/blob/master/SYNTAX.md.
func Get(data []byte, path string) gjson.Result {
	return gjson.GetBytes(data, path)
}

// Set replaces the value after searching the specified path in given JSON data,
// or deletes the value if given nil `value`,
// the path rule is documented on https://github.com/tidwall/gjson/blob/master/SYNTAX.md.
func Set(data []byte, path string, value []byte) ([]byte, error) {
	if value == nil {
		return sjson.DeleteBytes(data, path)
	}

	return sjson.SetRawBytes(data, path, value)
}
