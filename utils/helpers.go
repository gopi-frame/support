package utils

// With returns the given value, passed through the given callbacks.
func With[T any](value T, callbacks ...func(value T) T) T {
	for _, callback := range callbacks {
		value = callback(value)
	}
	return value
}

// Through call the given callbacks with the given value then return the value.
func Through[T any](value T, callbacks ...func(value T)) T {
	for _, callback := range callbacks {
		callback(value)
	}
	return value
}
