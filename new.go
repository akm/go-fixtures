package fixtures

// New creates a new instance of type T and applies the provided options to it.
func New[T any](opts ...func(*T)) *T {
	var rr T
	r := &rr
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewWithDefaults creates a new instance of type T, applies the default options first,
// and then applies the provided options to it.
func NewWithDefaults[T any](opts []func(*T), defaultOpts ...func(*T)) *T {
	options := append(defaultOpts, opts...)
	return New[T](options...)
}
