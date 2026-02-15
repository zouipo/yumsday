package mappers

// MapList converts a slice of one type to a slice of another using the provided mapping function.
func MapList[S any, D any](source []S, mapFunc func(*S) D) []D {
	result := make([]D, len(source))

	for i := range source {
		result[i] = mapFunc(&source[i])
	}
	return result
}
