package main

func Map[T, U any](arr []T, mapper func(val T) U) []U {
	u := make([]U, len(arr))

	for i, t := range arr {
		u[i] = mapper(t)
	}

	return u
}

func Filter[T any](arr []T, filter func(val T) bool) []T {
	res := make([]T, 0)

	for _, val := range arr {
		if filter(val) {
			res = append(res, val)
		}
	}

	return res
}
