package utils

func Filter[T any](ts []T, f func(T) bool) []T {
	us := make([]T, 0)
	for _, t := range ts {
		if f(t) {
			us = append(us, t)
		}
	}
	return us
}
