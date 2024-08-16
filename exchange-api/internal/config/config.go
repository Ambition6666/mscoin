package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ExchangeRPC zrpc.RpcClientConf
	CacheRedis  cache.CacheConf
	JWT         JWTConf
}

type JWTConf struct {
	AccessExpire int64
	AccessSecret string
}
