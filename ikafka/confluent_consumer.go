package ikafka

import (
	"strings"

	"github.com/airunny/wiki-go-tools/recovery"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kratos/kratos/v2/log"
)

type ConfluentConsumer struct {
	close    chan struct{}
	consumer *kafka.Consumer
	topics   []string
}

type ConfluentDo func(c *kafka.Consumer, msg *kafka.Message)

func NewConfluentConsumer(cfg *ConsumerConfig) (*ConfluentConsumer, error) {
	configMap := kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(cfg.Brokers, ","),
		"group.id":           cfg.GroupId,
		"enable.auto.commit": false,
	}

	if cfg.AutoCommit {
		configMap["enable.auto.commit"] = true
	}

	if cfg.OffsetOldest {
		configMap["auto.offset.reset"] = "earliest"
	}

	c, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(cfg.Topics, nil)
	if err != nil {
		return nil, err
	}

	out := &ConfluentConsumer{
		close:    make(chan struct{}),
		consumer: c,
		topics:   cfg.Topics,
	}
	out.log()
	return out, nil
}

func (s *ConfluentConsumer) Start(do ConfluentDo) {
	go func() {
		recovery.CatchGoroutinePanic()
		for {
			select {
			case <-s.close:
				log.Infof("ConfluentConsumer[%s] Stop.", s.topics)
				return
			default:
				ev := s.consumer.Poll(5000)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {
				case *kafka.Message:
					s.consumer.StoreMessage(e)
					do(s.consumer, e)
				case kafka.Error:
					log.Errorf("ConfluentConsumer[%s] Err:%v", s.topics, e)
					if e.Code() == kafka.ErrAllBrokersDown {
						return
					}
				}
			}
		}
	}()
}

func (s *ConfluentConsumer) Close() error {
	close(s.close)
	return s.consumer.Close()
}

func (s *ConfluentConsumer) log() {
	go func() {
		defer recovery.CatchGoroutinePanic()
		for l := range s.consumer.Logs() {
			switch l.Level {
			case 1: // DEBUG
				log.Debugf("kafka consumer <%v> msg:%s", s.topics, l.String())
			case 2: // INFO
				log.Infof("kafka consumer <%v> msg:%s", s.topics, l.String())
			case 3: // NOTICE
				log.Infof("kafka consumer <%v> msg:%s", s.topics, l.String())
			case 4: // WARNING
				log.Infof("kafka consumer <%v> msg:%s", s.topics, l.String())
			case 5: // ERROR
				log.Errorf("kafka consumer <%v> msg:%s", s.topics, l.String())
			default:
				log.Errorf("kafka consumer <%v> msg:%s", s.topics, l.String())
			}
		}
	}()
}
