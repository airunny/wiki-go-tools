package ikafka

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/airunny/wiki-go-tools/recovery"
	"github.com/go-kratos/kratos/v2/log"
)

// 事务类目前没有做，如果需要找yann

type ProducerConfig struct {
	Brokers      []string `json:"brokers" yaml:"brokers"`
	Username     string   `json:"username" yaml:"username"`
	Password     string   `json:"password" yaml:"password"`
	EnableSASL   bool     `json:"enable_sasl" yaml:"enableSASL"`
	EnableTLS    bool     `json:"enable_tls" yaml:"enableTLS"`
	Version      string   `json:"version" yaml:"version"`
	Topic        string   `json:"topic" yaml:"topic"`
	Algorithm    string   `json:"algorithm"`
	RequiredAcks int      `json:"required_acks"`
}

type Producer struct {
	transactionIdGenerator int32
	producersLock          sync.Mutex
	producers              []sarama.AsyncProducer
	producerProvider       func() sarama.AsyncProducer
}

func NewProducer(cfg *ProducerConfig) (*Producer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("empty brokers %v", cfg.Brokers)
	}

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, err
	}

	conFunc := newProducerConfig(cfg, version)

	p := &Producer{}
	p.producerProvider = func() sarama.AsyncProducer {
		config := conFunc()
		suffix := p.transactionIdGenerator
		if config.Producer.Transaction.ID != "" {
			p.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}

		producer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
		if err != nil {
			panic(err)
		}

		p.monitorErrors(producer)
		p.monitorSuccess(producer)
		return producer
	}
	return p, nil
}

func newProducerConfig(cfg *ProducerConfig, version sarama.KafkaVersion) func() *sarama.Config {
	return func() *sarama.Config {
		conf := sarama.NewConfig()
		conf.Version = version
		conf.Producer.Retry.Max = 3
		conf.Producer.RequiredAcks = sarama.RequiredAcks(cfg.RequiredAcks) // 保证消息不丢失
		conf.Producer.Return.Successes = true
		//conf.Producer.Idempotent = true                             // 保证幂等
		//conf.Producer.Return.Errors = true                          // 发生错误时需要记录
		conf.Producer.Partitioner = sarama.NewRoundRobinPartitioner // 轮训的方式
		//conf.Producer.Transaction.Retry.Backoff = 10                // 事务失败尝试时间间隔
		//conf.Producer.Transaction.ID = "producer"                   // 生产者ID
		//conf.Net.MaxOpenRequests = 1
		conf.Net.SASL.Enable = cfg.EnableSASL
		conf.Net.TLS.Enable = cfg.EnableTLS
		conf.Producer.Compression = sarama.CompressionGZIP
		if cfg.Username != "" && cfg.Password != "" {
			conf.Net.SASL.User = cfg.Username
			conf.Net.SASL.Password = cfg.Password
			conf.Net.TLS.Config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		if cfg.Algorithm == "sha512" {
			conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		} else if cfg.Algorithm == "sha256" {
			conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		}
		return conf
	}
}

func (p *Producer) monitorSuccess(producer sarama.AsyncProducer) {
	go func() {
		defer recovery.CatchGoroutinePanic()
		for success := range producer.Successes() {
			_ = success
		}
	}()
}

func (p *Producer) monitorErrors(producer sarama.AsyncProducer) {
	go func() {
		defer recovery.CatchGoroutinePanic()
		for err := range producer.Errors() {
			value, _ := err.Msg.Value.Encode()
			log.Errorf("kafka produce <%s,%s> Err:%v", err.Msg.Topic, string(value), err)
		}
	}()
}

func (p *Producer) AsyncMessage(topic string, value sarama.Encoder, opts ...MessageOption) {
	o := defaultMessageOption()
	for _, opt := range opts {
		opt(o)
	}
	p.asyncMessage(value, topic, o)
}

func (p *Producer) asyncMessage(value sarama.Encoder, topic string, o *messageOption) {
	producer := p.Borrow()
	defer p.Release(producer)
	producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       o.Key,
		Value:     value,
		Headers:   o.Headers,
		Metadata:  o.Metadata,
		Offset:    o.Offset,
		Partition: o.Partition,
		Timestamp: time.Now(),
	}
}

func (p *Producer) Borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *Producer) Release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *Producer) Close() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		err := producer.Close()
		if err != nil {
			log.Errorf("close producer err:%v", err)
		}
	}
	p.producers = p.producers[:0]
}
