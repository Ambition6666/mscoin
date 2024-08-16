package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql       MysqlConfig
	CacheRedis  cache.CacheConf
	Captcha     CaptchaConf
	JWT         JWTConf
	MarketRPC   zrpc.RpcClientConf
	Kafka       KafkaConfig
	ExchangeRPC zrpc.RpcClientConf
}
type CaptchaConf struct {
	Vid string
	Key string
}

type MysqlConfig struct {
	DataSource string
}

type JWTConf struct {
	AccessExpire int64
	AccessSecret string
}

type KafkaConfig struct {
	Addr     string
	WriteCap int
	ReadCap  int
	Group    string `json:"group,optional"`
}
