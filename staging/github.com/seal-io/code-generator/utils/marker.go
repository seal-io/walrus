package utils

import "strings"

func ParseMarker(line string) map[string]string {
	rm := map[string]string{}

	h, r, ok := strings.Cut(line, "=")
	for ; ok; h, r, ok = strings.Cut(r, "=") {
		h = strings.TrimSpace(h)
		r = strings.TrimSpace(r)

		if r != "" {
			switch r[0] {
			case '[', '{':
				if i := indexJSON(r); -1 < i && i < len(r)-1 {
					rm[h] = r[:i+1]
					r = r[i+2:]
					continue
				}
			default:
				if i := strings.Index(r, ","); -1 < i && i < len(r) {
					rm[h] = r[:i]
					r = r[i+1:]
					continue
				}
			}
		}

		rm[h] = r
	}

	return rm
}

func indexJSON(r string) int {
	var c int

	for i := 0; i < len(r); i++ {
		switch r[i] {
		case '{', '[':
			c++
		case '}', ']':
			c--
			if c == 0 {
				return i
			}
		}
	}

	return -1
}
