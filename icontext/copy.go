package icontext

import (
	"context"
	"time"
)

func Detach(ctx context.Context) context.Context {
	return dcontext{ctx}
}

type dcontext struct {
	ctx context.Context
}

func (dcontext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (dcontext) Done() <-chan struct{}       { return nil }
func (dcontext) Err() error                  { return nil }
func (x dcontext) Value(key interface{}) interface{} {
	return x.ctx.Value(key)
}
