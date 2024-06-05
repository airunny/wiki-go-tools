package locker

import (
	"context"
	"errors"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const (
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

func NewLockerWithRedis(cli *redis.Client) (Locker, error) {
	if cli == nil {
		return nil, errors.New("empty cli")
	}

	return &redisLocker{
		redisCli: cli,
	}, nil
}

type redisLocker struct {
	redisCli *redis.Client
}

func (s *redisLocker) Lock(ctx context.Context, key string, expires, timeout time.Duration) (Release, error) {
	client := s.redisCli.WithContext(ctx)

	value := time.Now().UnixNano()
	ok, err := client.SetNX(ctx, key, value, expires).Result()
	if err != nil {
		return nil, err
	}

	if ok {
		return s.Release(key, value), nil
	}

	timer := time.After(timeout)
	for {
		select {
		case <-time.After(time.Millisecond * 10):
		case <-timer:
			return nil, ErrLockTimeout
		}

		value = time.Now().UnixNano()
		ok, err = client.SetNX(ctx, key, value, expires).Result()
		if err != nil {
			return nil, err
		}

		if ok {
			return s.Release(key, value), nil
		}
	}
}

func (s *redisLocker) TryLock(ctx context.Context, key string, expires time.Duration) (Release, error) {
	client := s.redisCli.WithContext(ctx)

	value := time.Now().UnixNano()
	ok, err := client.SetNX(ctx, key, value, expires).Result()
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrAlreadyLocked
	}

	return s.Release(key, value), nil
}

func (s *redisLocker) Release(key string, value int64) Release {
	return func() error {
		val, err := s.redisCli.Eval(context.Background(), delCommand, []string{key}, value).Int64()
		if err != nil {
			return err
		}
		if val == 0 {
			return errors.New("release lock failed")
		}
		return nil
	}
}
