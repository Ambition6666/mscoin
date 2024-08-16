package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UCenterRPC zrpc.RpcClientConf
	JWT        JWTConf
	MarketRPC  zrpc.RpcClientConf
}

type JWTConf struct {
	AccessExpire int64
	AccessSecret string
}
