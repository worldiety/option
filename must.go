package option

// MustZero panics if the given value is not equal to the generic zero. This
// is also useful to assert a nil error.
func MustZero[T comparable](value T) {
	var zero T
	if zero != value {
		panic(value)
	}
}

// Must asserts that the tupel of (T,error) does not contain an error and
// otherwise panics.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
