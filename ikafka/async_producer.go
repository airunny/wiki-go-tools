package ikafka

import (
	"crypto/tls"

	"github.com/IBM/sarama"
	"github.com/airunny/wiki-go-tools/recovery"
	"github.com/go-kratos/kratos/v2/log"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
}

func NewAsyncProducer(cfg *ProducerConfig) (*AsyncProducer, error) {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, err
	}

	conf := sarama.NewConfig()
	conf.Version = version
	conf.Producer.Retry.Max = 3
	conf.Producer.RequiredAcks = sarama.RequiredAcks(cfg.RequiredAcks) // 保证消息不丢失
	conf.Producer.Return.Successes = true
	conf.Producer.Partitioner = sarama.NewRoundRobinPartitioner // 轮训的方式
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

	err = conf.Validate()
	if err != nil {
		return nil, err
	}

	ret := &AsyncProducer{}

	ret.producer, err = sarama.NewAsyncProducer(cfg.Brokers, conf)
	if err != nil {
		return nil, err
	}

	ret.monitorErrors()
	ret.monitorSuccess()
	return ret, nil
}

func (p *AsyncProducer) monitorSuccess() {
	go func() {
		defer recovery.CatchGoroutinePanic()
		for success := range p.producer.Successes() {
			_ = success
		}
	}()
}

func (p *AsyncProducer) monitorErrors() {
	go func() {
		defer recovery.CatchGoroutinePanic()
		for err := range p.producer.Errors() {
			value, _ := err.Msg.Value.Encode()
			log.Errorf("kafka produce <%s,%s> Err:%v", err.Msg.Topic, string(value), err)
		}
	}()
}

func (p *AsyncProducer) AsyncMessage(topic string, val sarama.Encoder, opts ...MessageOption) {
	o := defaultMessageOption()
	for _, opt := range opts {
		opt(o)
	}

	p.producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       o.Key,
		Value:     val,
		Headers:   o.Headers,
		Metadata:  o.Metadata,
		Offset:    o.Offset,
		Partition: o.Partition,
	}
}

func (p *AsyncProducer) Close() {
	p.producer.AsyncClose()
}
