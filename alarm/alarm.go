package alarm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/airunny/wiki-go-tools/env"
	"github.com/airunny/wiki-go-tools/icontext"
	"github.com/airunny/wiki-go-tools/recovery"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-resty/resty/v2"
)

const (
	atUsers      = `<at user_id="all">所有人</at> `
	atUserLayout = `<at user_id="ou_%s">%s</at>`
)

type Config struct {
	Delay        time.Duration `json:"delay"`
	NoDelayCount int           `json:"no_delay_count"`
	Webhook      string        `json:"webhook"`
	ENV          string        `json:"env"`
	ServiceName  string        `json:"service_name"`
	Users        []string      `json:"users"`
}

func NewAlarm(c *Config) (*Alarm, error) {
	if c.Delay <= 0 {
		c.Delay = time.Second * 10
	}

	if c.NoDelayCount <= 0 {
		c.NoDelayCount = 3
	}

	if c.ENV == "" {
		c.ENV = env.Environment()
	}

	if c.ServiceName == "" {
		c.ServiceName = env.GetServiceName()
	}

	if c.Webhook == "" {
		return nil, errors.New("empty webhook")
	}

	atUser := atUsers
	if len(c.Users) > 0 {
		users := make([]string, 0, len(c.Users))
		for _, user := range c.Users {
			users = append(users, fmt.Sprintf(atUserLayout, user, user))
		}
		atUser = strings.Join(users, "")
	}

	out := &Alarm{
		httpClient: resty.New().
			AddRetryCondition(RetryCondition).
			OnBeforeRequest(BeforeRequest),
		cfg:      c,
		close:    make(chan struct{}, 1),
		message:  make([]*message, 0, c.NoDelayCount),
		sendTime: time.Now(),
		atUser:   atUser,
	}
	out.Start()
	return out, nil
}

type message struct {
	RequestId string
	Message   string
}

type Alarm struct {
	sync.Mutex
	httpClient   *resty.Client
	cfg          *Config
	close        chan struct{}
	message      []*message
	messageCount int
	sendTime     time.Time
	atUser       string
}

func (a *Alarm) Start() {
	go func() {
		defer recovery.CatchGoroutinePanic()
		timer := time.NewTicker(a.cfg.Delay)
		defer timer.Stop()

		getMessages := func() []*message {
			sendMessages := make([]*message, 0, len(a.message))
			for _, item := range a.message {
				sendMessages = append(sendMessages, &message{
					RequestId: item.RequestId,
					Message:   item.Message,
				})
			}
			a.messageCount = 0
			a.message = a.message[:0]
			return sendMessages
		}

		for {
			select {
			case <-timer.C:
				a.Lock()
				if len(a.message) > 0 {
					sendMessages := getMessages()
					a.Unlock()
					go func() {
						defer recovery.CatchGoroutinePanic()
						a.send(sendMessages)
					}()
				} else {
					a.Unlock()
				}
			case <-a.close:
				a.Lock()
				if len(a.message) <= 0 {
					a.Unlock()
				} else {
					sendMessages := getMessages()
					a.Unlock()
					a.send(sendMessages)
					return
				}
			}
		}
	}()
}

func (a *Alarm) Alarm(reqId string, msg string) {
	log.Context(icontext.WithRequestId(context.Background(), reqId)).Error(msg)
	a.Lock()
	if len(a.message) <= 0 {
		a.sendTime = time.Now()
	}

	a.message = append(a.message, &message{
		RequestId: reqId,
		Message:   msg,
	})
	a.messageCount += 1

	if a.messageCount < a.cfg.NoDelayCount {
		a.Unlock()
		return
	}

	sendMessages := make([]*message, 0, len(a.message))
	for _, item := range a.message {
		sendMessages = append(sendMessages, &message{
			RequestId: item.RequestId,
			Message:   item.Message,
		})
	}
	a.messageCount = 0
	a.message = a.message[:0]
	a.Unlock()

	go func() {
		defer recovery.CatchGoroutinePanic()
		a.send(sendMessages)
	}()
}

func (a *Alarm) SendMessage(msg string) {
	_, err := a.httpClient.R().
		SetBody(a.Body(msg)).
		Post(a.cfg.Webhook)
	if err != nil {
		log.Errorf("Alarm:%v", msg)
	}
}

func (a *Alarm) AlarmNow(reqId string, msg string, opts ...Option) {
	o := &option{}
	for _, opt := range opts {
		opt(o)
	}

	l := log.Context(icontext.WithRequestId(context.Background(), reqId))
	switch o.level {
	case Debug:
		l.Debug(msg)
	case Info:
		l.Info(msg)
	case Warn:
		l.Warn(msg)
	case Error:
		l.Error(msg)
	}

	title := fmt.Sprintf("[%s][%v]", a.cfg.ServiceName, strings.ToUpper(a.cfg.ENV))
	_, err := a.httpClient.R().
		SetBody(a.Body(fmt.Sprintf("%s\n%s\n%s", a.atUser, title, msg))).
		Post(a.cfg.Webhook)
	if err != nil {
		log.Errorf("Alarm:[%v]%v", title, msg)
	}
}

func (a *Alarm) send(in []*message) {
	if len(in) <= 0 {
		return
	}

	dur := time.Now().Sub(a.sendTime)
	a.sendTime = time.Now()

	title := ""
	if a.cfg.ServiceName != "" {
		title += fmt.Sprintf("[%s]", a.cfg.ServiceName)
	}

	if a.cfg.ENV != "" {
		title += fmt.Sprintf("[%s]", strings.ToUpper(a.cfg.ENV))
	}

	title += fmt.Sprintf("%v时间内有%v条报警消息", dur, len(in))
	var contents []string

	for _, item := range in {
		contents = append(contents, fmt.Sprintf("[%v]%v", item.RequestId, item.Message))
	}

	_, err := a.httpClient.R().
		SetBody(a.Body(fmt.Sprintf("%s\n%s\n%s", a.atUser, title, strings.Join(contents, "\n")))).
		Post(a.cfg.Webhook)
	if err != nil {
		log.Errorf("Alarm:[%v]%v", title, strings.Join(contents, "\n"))
	}
}

func (a *Alarm) Close() {
	close(a.close)
}

func (a *Alarm) Body(text string) *Body {
	return &Body{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: text,
		},
	}
}

type Body struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func BeforeRequest(_ *resty.Client, request *resty.Request) error {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	return nil
}

func RetryCondition(response *resty.Response, _ error) bool {
	return response.StatusCode() != http.StatusOK
}
