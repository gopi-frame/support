package utils

// Ptr returns the pointer to the specific value
func Ptr[T any](value T) *T {
	return &value
}
