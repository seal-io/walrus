package io

func IsCreateInputDisabled(as map[string]any) bool {
	return HasAnnotation(as, "CreateInputDisabled")
}

func IsUpdateInputDisabled(as map[string]any) bool {
	return HasAnnotation(as, "UpdateInputDisabled")
}

func IsOutputDisabled(as map[string]any) bool {
	return HasAnnotation(as, "OutputDisabled")
}

func IsInline(as map[string]any) bool {
	return HasAnnotation(as, "Inline")
}

func HasAnnotation(as map[string]any, k string) bool {
	// Get desired annotation from annotation map.
	av, ave := as[annotationName]
	if !ave || av == nil {
		return false
	}

	// Convert typed annotation.
	an, anok := av.(map[string]any)
	if !anok || an == nil {
		return false
	}

	// Get value via given key.
	kv, kve := an[k]
	if !kve || kv == nil {
		return false
	}

	// Convert typed value.
	v, ok := kv.(bool)

	return ok && v
}
