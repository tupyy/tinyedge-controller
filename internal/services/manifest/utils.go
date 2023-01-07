package manifest

// substract return all elements of a which are not found in b
func substract[T, R any, S func(elem T) string, V func(elem R) string](a []T, b []R, idFn S, idFn2 V) []T {
	m1 := make(map[string]T)
	m2 := make(map[string]R)

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
			id := idFn2(b[i])
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

func substract1[T any, S func(elem T) string](a []T, b []T, idFn S) []T {
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

func intersect[T, R any, S func(elem T) string, V func(elem R) string, W func(elem T, elemR R) bool](a []T, b []R, idFn S, idFn2 V, equal W) []T {
	m1 := make(map[string]T)
	m2 := make(map[string]R)

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
			id := idFn2(b[i])
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
