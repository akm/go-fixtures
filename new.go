package fixtures

func New[T any](opts ...func(*T)) *T {
	var rr T
	r := &rr
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func NewWithDefaults[T any](opts []func(*T), defaultOpts ...func(*T)) *T {
	options := append(defaultOpts, opts...)
	return New[T](options...)
}
