package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/client"
	"market-api/internal/config"
	"market-api/internal/database"
	"market-api/internal/processor"
	"market-api/internal/ws"
)

type ServiceContext struct {
	Config          config.Config
	ExchangeRateRPC client.ExchangeRate
	MarketRPC       client.Market
	Cache           cache.Cache
	Kafka           *database.KafkaClient
	Processor       processor.Processor
}

func NewServiceContext(c config.Config, server *ws.WebSocketServer) *ServiceContext {
	newRedis := c.CacheRedis[0].NewRedis()
	node := cache.NewNode(newRedis, nil, nil, nil)
	kcli := database.NewKafkaClient(c.Kafka)

	p := processor.NewDefaultProcessor(kcli)
	p.Init(processor.NewWebsocketHandler(server))

	return &ServiceContext{
		Config:          c,
		ExchangeRateRPC: client.NewExchangeRate(zrpc.MustNewClient(c.MarketRPC)),
		MarketRPC:       client.NewMarket(zrpc.MustNewClient(c.MarketRPC)),
		Cache:           node,
		Kafka:           kcli,
		Processor:       p,
	}
}
