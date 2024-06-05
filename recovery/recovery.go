package recovery

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

func CatchGoroutinePanic() {
	if err := recover(); err != nil {
		log.Errorf("Recovery:panic:%v", err)
	}
}

func CatchGoroutinePanicWithContext(ctx context.Context) {
	if err := recover(); err != nil {
		log.Context(ctx).Errorf("Recovery:panic:%v", err)
	}
}
