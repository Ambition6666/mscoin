package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql      MysqlConfig
	CacheRedis cache.CacheConf
	MarketRPC  zrpc.RpcClientConf
	UCenterRPC zrpc.RpcClientConf
	Kafka      KafkaConfig
}

type MysqlConfig struct {
	DataSource string
}
type KafkaConfig struct {
	Addr     string
	WriteCap int
	ReadCap  int
	Group    string `json:"group,optional"`
}
