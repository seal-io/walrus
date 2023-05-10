package maps

import (
	"reflect"
)

// RemoveNulls takes a map of type map[string]interface{} and removes all nil values from it.
func RemoveNulls(m map[string]interface{}) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		case map[string]interface{}:
			RemoveNulls(t)
		case []interface{}:
			for _, v := range t {
				if t, ok := v.(map[string]interface{}); ok {
					RemoveNulls(t)
				}
			}
		case []map[string]interface{}:
			for _, v := range t {
				RemoveNulls(v)
			}
		}
	}
}

func RemoveNullsCopy(m map[string]interface{}) map[string]interface{} {
	newMap := CopyMap(m)
	RemoveNulls(newMap)

	return newMap
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// GetString gets a string value by key from a map of type map[string]interface{}.
func GetString(m map[string]interface{}, key string) string {
	v, exist := m[key]
	if !exist {
		return ""
	}
	vs, ok := v.(string)
	if !ok {
		return ""
	}
	return vs
}
