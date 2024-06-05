package iredis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/go-kratos/kratos/v2/log" // nolint
	redis "github.com/go-redis/redis/v8"
)

type Config struct {
	MasterName      string      `json:"master_name" yaml:"masterName" env:"DP_REDIS_DB_MASTER_NAME"`
	SentinelAddress []string    `json:"sentinel_address" yaml:"sentinelAddress" env:"DP_REDIS_DB_SENTINEL_ADDRESS"`
	Address         string      `json:"address" yaml:"address" env:"DP_REDIS_DB_SOURCE"`
	Password        string      `json:"password" yaml:"password" env:"DP_REDIS_DB_PASSWORD"`
	DB              int         `json:"db" yaml:"db" env:"DP_REDIS_DB_NUMBER"`
	MaxIdle         int         `json:"max_idle" yaml:"max_idle" env:"DP_REDIS_DB_MAX_IDLE"`
	CaPath          string      `json:"ca_path" yaml:"caPath"`
	TLSConfig       *tls.Config `json:"tls_config" yaml:"tlsConfig"`
	Debug           bool        `json:"debug" yaml:"debug"`
}

func NewClient(cfg *Config, l log.Logger) (*redis.Client, error) {
	if cfg.MaxIdle <= 0 {
		cfg.MaxIdle = 30
	}

	if len(cfg.CaPath) > 0 {
		sslCfg, err := GenTLSConfig(cfg.CaPath)
		if err != nil {
			return nil, err
		}
		cfg.TLSConfig = sslCfg
	}

	var redisClient *redis.Client
	if len(cfg.SentinelAddress) != 0 {
		redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			SentinelAddrs: cfg.SentinelAddress,
			Password:      cfg.Password,
			DB:            cfg.DB,
			PoolSize:      cfg.MaxIdle,
			TLSConfig:     cfg.TLSConfig,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:      cfg.Address,
			Password:  cfg.Password,
			DB:        cfg.DB,
			PoolSize:  cfg.MaxIdle,
			TLSConfig: cfg.TLSConfig,
		})
	}

	ll := newLogger(l, cfg.Debug)
	redis.SetLogger(ll)
	redisClient.AddHook(newHook(ll))
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}

func GenTLSConfig(caCertFile string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}, nil
}
