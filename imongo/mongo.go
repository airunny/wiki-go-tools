package imongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Config struct {
		URL         string      `json:"url" yaml:"url" env:"DP_MONGO_DB_URL"`
		MaxOpen     uint64      `json:"maxOpen" yaml:"maxOpen" env:"DP_MONGO_MAX_OPEN"`
		MaxIdle     int         `json:"maxIdle" yaml:"maxIdle" env:"DP_MONGO_MAX_IDLE"`
		MaxLifeTime int         `json:"maxLifeTime" yaml:"maxLifeTime" env:"DP_MONGO_MAX_LIFE_TIME"`
		CaPath      string      `json:"ca_path" yaml:"caPath"`
		TlsConfig   *tls.Config `json:"tls_config" yaml:"tlsConfig"`
		Debug       bool        `json:"debug" yaml:"debug"`
	}

	Client struct {
		*mongo.Client
		debug bool
	}
)

func NewClient(cfg *Config) (*Client, error) {
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()

	opts := options.Client()
	opts.ApplyURI(cfg.URL)

	if cfg.MaxOpen > 0 {
		opts.SetMaxPoolSize(cfg.MaxOpen)
	}

	if cfg.MaxLifeTime > 0 {
		opts.SetMaxConnIdleTime(time.Second * time.Duration(cfg.MaxLifeTime))
	}

	if len(cfg.CaPath) > 0 {
		sslCfg, err := GenTLSConfig(cfg.CaPath)
		if err != nil {
			return nil, err
		}
		opts.SetTLSConfig(sslCfg)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	newCtx, newCancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer newCancel()

	err = client.Ping(newCtx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	out := &Client{
		Client: client,
		debug:  cfg.Debug,
	}
	defaultClient = out

	return out, nil
}

// IDatabase 增加了比较全的链路追踪和日志信息
func (c *Client) IDatabase(name string, opts ...*options.DatabaseOptions) *Database {
	return &Database{
		Database: c.Database(name, opts...),
		debug:    c.debug,
	}
}

func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	return c.Client.Disconnect(context.Background())
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
