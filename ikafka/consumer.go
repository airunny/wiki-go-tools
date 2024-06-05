package ikafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

type ConsumerConfig struct {
	Brokers           []string      `json:"brokers" yaml:"brokers"`
	Topics            []string      `json:"topics" yaml:"topics"`
	GroupId           string        `json:"group_id" yaml:"group_id"`
	Username          string        `json:"username" yaml:"username"`
	Password          string        `json:"password" yaml:"password"`
	OffsetOldest      bool          `json:"offset_oldest" yaml:"offset_oldest"`
	EnableSASL        bool          `json:"enable_sasl" yaml:"enable_sasl"`
	EnableTLS         bool          `json:"enable_tls" yaml:"enable_tls"`
	Version           string        `json:"version" yaml:"version"`
	AutoCommit        bool          `json:"auto_commit" yaml:"auto_commit"`
	MaxProcessingTime time.Duration `json:"max_processing_time" yaml:"max_processing_time"`
}

type Do func(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage)

type Consumer struct {
	cfg    *ConsumerConfig
	ready  chan bool
	wg     *sync.WaitGroup
	client sarama.ConsumerGroup
	ctx    context.Context
	cancel context.CancelFunc
	do     Do
	atomic atomic.Int32
}

func NewConsumer(cfg *ConsumerConfig) (*Consumer, error) {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = version
	config.Net.SASL.Enable = cfg.EnableSASL
	config.Net.TLS.Enable = cfg.EnableTLS
	//config.Consumer.IsolationLevel = sarama.ReadCommitted

	if cfg.Username != "" {
		config.Net.SASL.User = cfg.Username
		config.Net.SASL.Password = cfg.Password
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if cfg.OffsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	if cfg.MaxProcessingTime > 0 {
		config.Consumer.MaxProcessingTime = cfg.MaxProcessingTime
	}

	config.Consumer.Offsets.AutoCommit.Enable = false
	//config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	//sarama.Logger = &Log{}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupId, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		cfg:    cfg,
		ready:  make(chan bool),
		wg:     &sync.WaitGroup{},
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (s *Consumer) Start(do Do) error {
	if value := s.atomic.Load(); value == 1 {
		return errors.New("start func can be executed only once")
	}

	s.atomic.Add(1)

	if do == nil {
		return errors.New("empty do func")
	}

	s.do = do
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			if err := s.client.Consume(s.ctx, s.cfg.Topics, s); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				panic(fmt.Sprintf("Consumer Consume Err:%v", err))
			}
			log.Infof("%s consumer restart", s.cfg.Topics)

			if s.ctx.Err() != nil {
				return
			}
			s.ready = make(chan bool)
		}
	}()

	go func() {
		<-s.ready
		log.Info("consumer running")
	}()
	return nil
}

func (s *Consumer) Close() error {
	s.cancel()
	s.wg.Wait()
	return s.client.Close()
}

func (s *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(s.ready)
	return nil
}

func (s *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Infof("%v message channel was closed", s.cfg.Topics)
				return nil
			}
			session.MarkMessage(message, "")
			s.do(session, message)
		case <-session.Context().Done():
			return nil
		}
	}
}
