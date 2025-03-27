package excel

type options struct {
	sheetName string
}

type Option func(*options)

func WithSheetName(sheetName string) Option {
	return func(o *options) {
		o.sheetName = sheetName
	}
}
