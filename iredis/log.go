package iredis

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type logger struct {
	log   *log.Helper
	debug bool
}

func newLogger(l log.Logger, debug bool) *logger {
	return &logger{
		log:   log.NewHelper(l),
		debug: debug,
	}
}

func (l *logger) Printf(ctx context.Context, format string, v ...interface{}) {
	if l.debug {
		ll := l.log.WithContext(ctx)
		ll.Infof(format, v...)
	}
}
