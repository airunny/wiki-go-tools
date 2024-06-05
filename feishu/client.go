package feishu

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/airunny/wiki-go-tools/locker"
	redis "github.com/go-redis/redis/v8"
	resty "github.com/go-resty/resty/v2"
	goCache "github.com/liyanbing/go-cache"
	redisCache "github.com/liyanbing/go-cache/cacher/redis"
)

// Config 目前只支持一对app_id跟app_secret
type Config struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type FSClient struct {
	Config     *Config
	httpClient *resty.Client
	pool       sync.Pool
	cache      goCache.Cache
	locker     locker.Locker
}

func New(c *Config, rc *redis.Client) (*FSClient, error) {
	if c == nil {
		return nil, fmt.Errorf("empty config")
	}

	if c.AppID == "" {
		return nil, fmt.Errorf("empty app_id")
	}

	if c.AppSecret == "" {
		return nil, fmt.Errorf("empty app_secret")
	}

	if rc == nil {
		return nil, fmt.Errorf("empty redis client")
	}

	lk, err := locker.NewLockerWithRedis(rc)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}

	return &FSClient{
		Config: c,
		httpClient: resty.New().
			OnBeforeRequest(beforeRequest).
			OnAfterResponse(afterResponse).
			SetTransport(tr),
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
		cache:  redisCache.NewRedisCache(rc),
		locker: lk,
	}, nil
}

func (c *FSClient) getBuffer() *bytes.Buffer {
	buf := c.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (c *FSClient) putBuffer(buf *bytes.Buffer) {
	c.pool.Put(buf)
}

func (c *FSClient) httpGet(ctx context.Context, in *doRequest) error {
	req := c.httpClient.R()
	if in.authorization != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %v", in.authorization))
	}

	_, err := req.
		SetContext(ctx).
		SetResult(in.out).
		Get(in.domain)
	return err
}

func (c *FSClient) httpPost(ctx context.Context, in *doRequest) error {
	buf := c.getBuffer()
	defer c.putBuffer(buf)

	err := json.NewEncoder(buf).Encode(in.req)
	if err != nil {
		return err
	}

	req := c.httpClient.R()
	if in.authorization != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %v", in.authorization))
	}

	_, err = req.
		SetContext(ctx).
		SetBody(buf).
		SetResult(in.out).
		Post(in.domain)
	if err != nil {
		return err
	}
	return nil
}

func beforeRequest(_ *resty.Client, request *resty.Request) error {
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	request.Header.Set("Accept", "application/json;charset=utf-8")
	return nil
}

func afterResponse(_ *resty.Client, response *resty.Response) error {
	fmt.Println("原始数据：", string(response.Body()))
	if response.StatusCode() <= 199 || response.StatusCode() >= 300 {
		if response.Request.Result != nil {
			err := json.Unmarshal(response.Body(), response.Request.Result)
			if err != nil {
				return err
			}
		}
	}

	ret := response.Result()
	if check, ok := ret.(CheckResponse); ok {
		return check.Check()
	}
	return nil
}

type CheckResponse interface {
	Check() error
}

type doRequest struct {
	domain        string
	req           interface{}
	out           interface{}
	authorization string
}
