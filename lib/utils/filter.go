package utils

func Filter[T any](inputArray []T, predicate func(T) bool) []T {
	outputArray := make([]T, 0)
	for _, t := range inputArray {
		if predicate(t) {
			outputArray = append(outputArray, t)
		}
	}
	return outputArray
}
