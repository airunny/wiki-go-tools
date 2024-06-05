package urlformat

type options struct {
	template string
}

type Option func(o *options)

func WithTemplate(in string) Option {
	return func(o *options) {
		o.template = in
	}
}
