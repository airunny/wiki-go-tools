package ikafka

import (
	"fmt"
	"github.com/IBM/sarama"
)

func NewSyncProducer(cfg *ProducerConfig) (sarama.SyncProducer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("empty brokers %v", cfg.Brokers)
	}

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, err
	}

	conFunc := newProducerConfig(cfg, version)
	return sarama.NewSyncProducer(cfg.Brokers, conFunc())
}
