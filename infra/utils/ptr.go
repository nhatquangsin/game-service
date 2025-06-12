package utils

// Of returns a pointer to the given value
func Of[T any](x T) *T {
	return &x
}
