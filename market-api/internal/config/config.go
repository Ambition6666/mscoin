package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	MarketRPC  zrpc.RpcClientConf
	CacheRedis cache.CacheConf
	Kafka      KafkaConfig
}

type KafkaConfig struct {
	Addr     string
	WriteCap int
	ReadCap  int
	Group    string `json:"group,optional"`
}
