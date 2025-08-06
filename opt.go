package envmap

type Opt func(o *config)

type config struct {
	Prefix string
}

func (o *config) apply(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithPrefix(prefix string) Opt {
	return func(o *config) {
		o.Prefix = prefix
	}
}
