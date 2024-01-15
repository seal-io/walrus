package maps

import (
	"reflect"
)

// RemoveNulls takes a map of type map[string]any and removes all nil values from it.
func RemoveNulls(m map[string]any) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}

		switch t := v.Interface().(type) {
		case map[string]any:
			RemoveNulls(t)
		case []any:
			for _, v := range t {
				if t, ok := v.(map[string]any); ok {
					RemoveNulls(t)
				}
			}
		case []map[string]any:
			for _, v := range t {
				RemoveNulls(v)
			}
		}
	}
}

func RemoveNullsCopy(m map[string]any) map[string]any {
	newMap := CopyMap(m)
	RemoveNulls(newMap)

	return newMap
}

func CopyMap(m map[string]any) map[string]any {
	cp := make(map[string]any)

	for k, v := range m {
		vm, ok := v.(map[string]any)
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// GetString gets a string value by key from a map of type map[string]any.
func GetString(m map[string]any, key string) string {
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
