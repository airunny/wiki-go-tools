package oss

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denverdino/aliyungo/sts"
)

const (
	defaultPolicy             = `{"Statement":[{"Action":["oss.PutObject"],"Effect":"Allow","Resource":["acs:oss:*:1624820479181668:wkstock-dev/*","acs:oss:*:1624820479181668:wkstock-pro/*"]}],"Version":"1"}` // nolint
	dateLayout                = "20060102"
	defaultSTSDurationSeconds = 3600
	defaultRoleSessionName    = "wkstock-oss-sts"
)

type Config struct {
	Bucket             string
	Endpoint           string
	AccessKeyID        string
	AccessKeySecret    string
	RoleArn            string
	Policy             string
	CDNHost            string
	RoleSessionName    string
	STSDurationSeconds int
}

func NewClient(c *Config) (*Client, error) {
	if c.Bucket == "" {
		return nil, errors.New("empty buckets")
	}

	if c.Endpoint == "" {
		return nil, errors.New("empty endpoint")
	}

	if c.AccessKeyID == "" {
		return nil, errors.New("empty access_key_id")
	}

	if c.AccessKeySecret == "" {
		return nil, errors.New("empty access_key_secret")
	}

	if c.Policy == "" {
		c.Policy = defaultPolicy
	}

	if c.STSDurationSeconds <= 0 {
		c.STSDurationSeconds = defaultSTSDurationSeconds
	}

	if c.RoleSessionName == "" {
		c.RoleSessionName = defaultRoleSessionName
	}

	client, err := oss.New(c.Endpoint, c.AccessKeyID, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(c.Bucket)
	if err != nil {
		return nil, err
	}

	stsClient := sts.NewClient(c.AccessKeyID, c.AccessKeySecret)
	stsClient.SetTransport(Transport())

	return &Client{
		cfg:       c,
		Client:    client,
		bucket:    bucket,
		stsClient: stsClient,
	}, nil
}

type Client struct {
	cfg *Config
	*oss.Client
	bucket    *oss.Bucket
	stsClient *sts.STSClient
}

func Transport() http.RoundTripper {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	handshakeTimeoutStr, ok := os.LookupEnv("TLSHandshakeTimeout")
	if ok {
		handshakeTimeout, err := strconv.Atoi(handshakeTimeoutStr)
		if err != nil {
			log.Printf("Get TLSHandshakeTimeout from env error: %v.", err)
		} else {
			t.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second
		}
	}

	responseHeaderTimeoutStr, ok := os.LookupEnv("ResponseHeaderTimeout")
	if ok {
		responseHeaderTimeout, err := strconv.Atoi(responseHeaderTimeoutStr)
		if err != nil {
			log.Printf("Get ResponseHeaderTimeout from env error: %v.", err)
		} else {
			t.ResponseHeaderTimeout = time.Duration(responseHeaderTimeout) * time.Second
		}
	}

	expectContinueTimeoutStr, ok := os.LookupEnv("ExpectContinueTimeout")
	if ok {
		expectContinueTimeout, err := strconv.Atoi(expectContinueTimeoutStr)
		if err != nil {
			log.Printf("Get ExpectContinueTimeout from env error: %v.", err)
		} else {
			t.ExpectContinueTimeout = time.Duration(expectContinueTimeout) * time.Second
		}
	}

	idleConnTimeoutStr, ok := os.LookupEnv("IdleConnTimeout")
	if ok {
		idleConnTimeout, err := strconv.Atoi(idleConnTimeoutStr)
		if err != nil {
			log.Printf("Get IdleConnTimeout from env error: %v.", err)
		} else {
			t.IdleConnTimeout = time.Duration(idleConnTimeout) * time.Second
		}
	}

	return t
}
func (c *Client) Upload(fileName, contentType string, reader io.Reader) (domain, url string, err error) {
	var (
		folderName     = time.Now().Format(dateLayout)
		yunFileTmpPath = folderName + "/" + fileName
	)

	err = c.bucket.PutObject(yunFileTmpPath, reader, oss.ContentType(contentType))
	if err != nil {
		return c.cfg.CDNHost, url, err
	}

	return c.cfg.CDNHost, "/" + yunFileTmpPath, nil
}

func (c *Client) Config() *Config {
	return c.cfg
}

func (c *Client) GetSTS() (*sts.AssumedRoleUserCredentials, error) {
	resp, err := c.stsClient.AssumeRole(sts.AssumeRoleRequest{
		RoleArn:         c.cfg.RoleArn,
		RoleSessionName: c.cfg.RoleSessionName,
		DurationSeconds: c.cfg.STSDurationSeconds,
		// Policy:          c.cfg.Policy,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Credentials, nil
}
