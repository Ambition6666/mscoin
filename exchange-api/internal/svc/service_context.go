package svc

import (
	"exchange-api/internal/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/exchange/client"
)

type ServiceContext struct {
	Config   config.Config
	OrderRPC client.Order
	Cache    cache.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	newRedis := c.CacheRedis[0].NewRedis()
	node := cache.NewNode(newRedis, nil, nil, nil)

	return &ServiceContext{
		Config:   c,
		OrderRPC: client.NewOrder(zrpc.MustNewClient(c.ExchangeRPC)),
		Cache:    node,
	}
}
