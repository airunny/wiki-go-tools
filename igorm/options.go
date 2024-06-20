package igorm

import (
	"gorm.io/gorm"
)

type Options struct {
	tx *gorm.DB
}

func (o *Options) Session() *gorm.DB {
	return o.tx
}

type Option func(o *Options)

func WithTransaction(tx *gorm.DB) Option {
	return func(o *Options) {
		o.tx = tx
	}
}

func NewOptions(tx *gorm.DB, opts ...Option) *Options {
	o := &Options{
		tx: tx,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
