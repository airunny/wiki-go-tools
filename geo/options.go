package geo

type options struct {
	languageCode string
}

type Option func(o *options)

func WithLanguage(languageCode string) Option {
	return func(o *options) {
		o.languageCode = languageCode
	}
}

func getOptions(opts ...Option) *options {
	o := &options{
		languageCode: "en",
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
