package lo

func Filter[T any](collection []T, predicate func(item T) bool) []T {
	result := make([]T, 0, len(collection))

	for _, item := range collection {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}
