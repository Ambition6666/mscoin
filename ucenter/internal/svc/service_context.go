package svc

import (
	"common/msdb"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	eclient "grpc-common/exchange/client"
	mclient "grpc-common/market/client"
	"ucenter/consumer"
	"ucenter/database"
	"ucenter/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	Cache     cache.Cache
	Db        *msdb.MsDB
	MarketRPC mclient.Market
	Kcli      *database.KafkaClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	newRedis := c.CacheRedis[0].NewRedis()
	db := database.ConnMysql(c.Mysql.DataSource)
	node := cache.NewNode(newRedis, nil, nil, nil)
	cli := database.NewKafkaClient(c.Kafka)
	cli.StartRead("add-exchange-asset")
	order := eclient.NewOrder(zrpc.MustNewClient(c.ExchangeRPC))
	go consumer.ExchangeOrderAddConsumer(cli, db, order, newRedis)
	cli = cli.StartReadNew("exchange_order_complete_update_success")
	go consumer.ExchangeOrderComplete(newRedis, cli, db)
	cli = cli.StartReadNew("BtcTransactionTopic")
	go consumer.BitCoinTransaction(node, cli, db)
	withdrawClient := cli.StartReadNew("withdraw")
	go consumer.WithdrawConsumer(withdrawClient, db)

	return &ServiceContext{
		Config:    c,
		Cache:     node,
		Db:        db,
		MarketRPC: mclient.NewMarket(zrpc.MustNewClient(c.MarketRPC)),
		Kcli:      database.NewKafkaClient(c.Kafka),
	}
}
