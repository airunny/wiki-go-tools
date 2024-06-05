package ikafka

import "github.com/IBM/sarama"

type messageOption struct {
	Key       sarama.Encoder
	Headers   []sarama.RecordHeader
	Metadata  interface{}
	Offset    int64
	Partition int32
}

func defaultMessageOption() *messageOption {
	return &messageOption{
		Key:       nil,
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
	}
}

type MessageOption func(o *messageOption)

func WithMessageKey(in sarama.Encoder) MessageOption {
	return func(o *messageOption) {
		o.Key = in
	}
}

func WithMessageHeaders(in []sarama.RecordHeader) MessageOption {
	return func(o *messageOption) {
		o.Headers = in
	}
}

func WithMessageMetadata(in interface{}) MessageOption {
	return func(o *messageOption) {
		o.Metadata = in
	}
}

func WithMessageOffset(in int64) MessageOption {
	return func(o *messageOption) {
		o.Offset = in
	}
}

func WithMessagePartition(in int32) MessageOption {
	return func(o *messageOption) {
		o.Partition = in
	}
}

// =============
