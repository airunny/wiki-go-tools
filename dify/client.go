package dify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	ResourceNotFound = errors.New("resource not found")
)

type Client struct {
	httpClient *resty.Client
	url        string
}

func NewClient(url, apiKey string) (*Client, error) {
	if url == "" {
		return nil, errors.New("empty dify url")
	}

	return &Client{
		httpClient: resty.New().
			SetHeader("Connection", "keep-alive").
			SetTransport(&http.Transport{
				MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
				IdleConnTimeout:     10 * time.Minute, // 空闲连接超时时间
			}).
			SetTimeout(10 * time.Minute).
			OnBeforeRequest(BeforeRequestWrap(Authorization(apiKey))).
			AddRetryCondition(RetryCondition).
			SetRetryCount(3).
			OnAfterResponse(AfterResponse),
		url: url,
	}, nil
}

func Authorization(apiKey string) string {
	return fmt.Sprintf("Bearer %s", apiKey)
}

func BeforeRequestWrap(auth string) resty.RequestMiddleware {
	return func(_ *resty.Client, request *resty.Request) error {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		if request.Header.Get("Authorization") == "" {
			request.Header.Set("Authorization", auth)
		}
		return nil
	}
}

func RetryCondition(response *resty.Response, _ error) bool {
	return response.StatusCode() >= http.StatusInternalServerError
}

func AfterResponse(_ *resty.Client, response *resty.Response) error {
	if response.StatusCode() == http.StatusNotFound {
		return ResourceNotFound
	}

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
