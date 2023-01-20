package reference

// substract return all elements of a which are not found in b
func substract[T any, S func(elem T) string](a []T, b []T, idFn S) []T {
	m1 := make(map[string]T)
	m2 := make(map[string]T)

	limit := len(a)
	if limit < len(b) {
		limit = len(b)
	}

	for i := 0; i < limit; i++ {
		if i < len(a) {
			id := idFn(a[i])
			m1[id] = a[i]
		}

		if i < len(b) {
			id := idFn(b[i])
			m2[id] = b[i]
		}
	}

	res := make([]T, 0, len(a))
	for id, v := range m1 {
		if _, found := m2[id]; !found {
			res = append(res, v)
		}
	}

	return res
}

func intersect[T any, S func(elem T) string, W func(e1, e2 T) bool](a []T, b []T, idFn S, equal W) []T {
	m1 := make(map[string]T)
	m2 := make(map[string]T)

	limit := len(a)
	if limit < len(b) {
		limit = len(b)
	}

	for i := 0; i < limit; i++ {
		if i < len(a) {
			id := idFn(a[i])
			m1[id] = a[i]
		}

		if i < len(b) {
			id := idFn(b[i])
			m2[id] = b[i]
		}
	}

	res := make([]T, 0, len(a))
	for id, v := range m1 {
		if vv, found := m2[id]; found && equal(v, vv) {
			res = append(res, v)
		}
	}

	return res
}

func contains(arr []string, id string) bool {
	for _, a := range arr {
		if a == id {
			return true
		}
	}
	return false
}
