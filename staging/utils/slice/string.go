package slice

func ContainsAny[T comparable](slice []T, items ...T) bool {
	for i := 0; i < len(items); i++ {
		for j := 0; j < len(slice); j++ {
			if slice[j] == items[i] {
				return true
			}
		}
	}
	return false
}

func ContainsAll[T comparable](slice []T, items ...T) bool {
	for i := 0; i < len(items); i++ {
		if !ContainsAny(slice, items[i]) {
			return false
		}
	}
	return true
}
