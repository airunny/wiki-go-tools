package config

import (
	"github.com/airunny/wiki-go-tools/ikafka"
)

type DataConfig_Database struct {
	Driver             string `protobuf:"bytes,1,opt,name=driver,json=driver,proto3" json:"driver"`
	Source             string `protobuf:"bytes,2,opt,name=source,json=source,proto3" json:"source"`
	Level              int32  `protobuf:"varint,3,opt,name=level,json=level,proto3" json:"level"`
	MaxOpen            int32  `protobuf:"varint,4,opt,name=max_open,json=maxOpen,proto3" json:"max_open"`
	MaxIdle            int32  `protobuf:"varint,5,opt,name=max_idle,json=maxIdle,proto3" json:"max_idle"`
	MaxLifeTimeSeconds int32  `protobuf:"varint,6,opt,name=max_life_time_seconds,json=maxLifeTimeSeconds,proto3" json:"max_life_time_seconds"`
}

type DataConfig_Redis struct {
	Address             string `protobuf:"bytes,1,opt,name=address,json=address,proto3" json:"address"`
	Password            string `protobuf:"bytes,2,opt,name=password,json=password,proto3" json:"password"`
	Db                  int32  `protobuf:"varint,3,opt,name=db,json=db,proto3" json:"db"`
	MaxIdle             int32  `protobuf:"varint,4,opt,name=max_idle,json=maxIdle,proto3" json:"max_idle"`
	ReadTimeoutSeconds  int64  `protobuf:"varint,5,opt,name=read_timeout_seconds,json=readTimeoutSeconds,proto3" json:"read_timeout_seconds"`
	WriteTimeoutSeconds int64  `protobuf:"varint,6,opt,name=write_timeout_seconds,json=writeTimeoutSeconds,proto3" json:"write_timeout_seconds"`
}

type DataConfig struct {
	Database *DataConfig_Database `protobuf:"bytes,1,opt,name=database,json=database,proto3" json:"database"`
	Redis    *DataConfig_Redis    `protobuf:"bytes,2,opt,name=redis,json=redis,proto3" json:"redis"`
	Mongodb  *DataConfig_Database `protobuf:"bytes,3,opt,name=mongodb,json=mongodb,proto3" json:"mongodb"`
}

type ServerConfig_HTTP struct {
	Network        string `protobuf:"bytes,1,opt,name=network,json=network,proto3" json:"network"`
	Addr           string `protobuf:"bytes,2,opt,name=addr,json=addr,proto3" json:"addr"`
	TimeoutSeconds int64  `protobuf:"varint,3,opt,name=timeout_seconds,json=timeoutSeconds,proto3" json:"timeout_seconds"`
}

type ServerConfig_GRPC struct {
	Network        string `protobuf:"bytes,1,opt,name=network,json=network,proto3" json:"network"`
	Addr           string `protobuf:"bytes,2,opt,name=addr,json=addr,proto3" json:"addr"`
	TimeoutSeconds int64  `protobuf:"varint,3,opt,name=timeout_seconds,json=timeoutSeconds,proto3" json:"timeout_seconds"`
}

type ServerConfig struct {
	Http      *ServerConfig_HTTP `protobuf:"bytes,1,opt,name=http,json=http,proto3" json:"http"`                  // http信息
	Grpc      *ServerConfig_GRPC `protobuf:"bytes,2,opt,name=grpc,json=grpc,proto3" json:"grpc"`                  // grpc信息
	OssDomain string             `protobuf:"bytes,4,opt,name=oss_domain,json=ossDomain,proto3" json:"oss_domain"` // 对象存储domain
}

type Bootstrap struct {
	Server   *ServerConfig `json:"server"`
	Data     *DataConfig   `json:"data"`
	Business *Business     `json:"business"`
}

type Consumer struct {
	Brokers      []string `json:"brokers"`
	Topics       []string `json:"topics"`
	GroupId      string   `json:"group_id"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Version      string   `json:"version"`
	OffsetOldest bool     `json:"offset_oldest"`
}

type AWS struct {
	Id       string `json:"id"`
	Secret   string `json:"secret"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
	ApiKey   string `json:"api_key"`
	Limit    int64  `json:"limit"`
}

type Elastic struct {
	Source   string `json:"source"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Business struct {
	Aws                   *AWS                   `json:"aws"`
	Elastic               *Elastic               `json:"elastic"`
	UserIdKey             string                 `json:"user_id_key"`
	AppFlyerLogConsumer   *Consumer              `json:"app_flyer_log_consumer"`
	SearchRecordConsumer  *Consumer              `json:"search_record_consumer"`
	SearchCollectConsumer *Consumer              `json:"search_collect_consumer"`
	FXStatisticsConsumer  *Consumer              `json:"fx_statistics_consumer"`
	KAFKAProducer         *ikafka.ProducerConfig `json:"kafka_producer"`
}
