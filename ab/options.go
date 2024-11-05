package ab

type options struct {
	fixedId map[string]string
}

func newDefaultOptions() *options {
	return &options{
		fixedId: make(map[string]string),
	}
}

type Option func(o *options)

func WithFixedId(group string, ids ...string) Option {
	return func(o *options) {
		if o.fixedId == nil {
			o.fixedId = make(map[string]string, len(ids))
		}
		for _, id := range ids {
			o.fixedId[id] = group
		}
	}
}
