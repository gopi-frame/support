package utils

// Must return the value, but panics when err is not nil
func Must[T any](v1 T, err error) T {
	if err != nil {
		panic(err)
	}
	return v1
}

// Must2 return the values, but panics when err is not nil
func Must2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// Must3 return the values, but panics when err is not nil
func Must3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3
}

// Must4 return the values, but panics when err is not nil
func Must4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4
}
