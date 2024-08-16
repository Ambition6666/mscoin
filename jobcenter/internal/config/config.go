package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Okx        OkxConf
	Mongo      MongoConf
	Kafka      KafkaConfig
	CacheRedis redis.RedisConf
	UCenterRPC zrpc.RpcClientConf
}

type OkxConf struct {
	ApiKey    string
	SecretKey string
	Pass      string
	Host      string
	Proxy     string
}

type MongoConf struct {
	Url      string
	Username string
	Password string
	Database string
}

type KafkaConfig struct {
	Addr     string
	WriteCap int
	ReadCap  int
	Group    string `json:"group,optional"`
}
