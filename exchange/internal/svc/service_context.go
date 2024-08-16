package svc

import (
	"common/msdb"
	"exchange/database"
	"exchange/internal/config"
	"exchange/internal/consumer"
	"exchange/internal/processor"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	mclient "grpc-common/market/client"
	uclient "grpc-common/ucenter/client"
)

type ServiceContext struct {
	Config       config.Config
	Cache        cache.Cache
	Db           *msdb.MsDB
	MarketRPC    mclient.Market
	MemberRPC    uclient.Member
	AssetRPC     uclient.Asset
	Kafka        *database.KafkaClient
	TradeFactory *processor.CoinTradeFactory
}

func (c *ServiceContext) Init() {
	factory := processor.InitCoinTradeFactory()
	factory.Init(c.MarketRPC, c.Kafka, c.Db)
	c.TradeFactory = factory
	kafkaConsumer := consumer.NewKafkaConsumer(c.Config.Kafka, factory, c.Db)
	kafkaConsumer.Run()
}
func NewServiceContext(c config.Config) *ServiceContext {
	newRedis := c.CacheRedis[0].NewRedis()
	node := cache.NewNode(newRedis, nil, nil, nil)
	kcli := database.NewKafkaClient(c.Kafka)
	kcli.StartWrite()
	db := database.ConnMysql(c.Mysql.DataSource)
	s := &ServiceContext{
		Config:    c,
		Cache:     node,
		Db:        db,
		MarketRPC: mclient.NewMarket(zrpc.MustNewClient(c.MarketRPC)),
		MemberRPC: uclient.NewMember(zrpc.MustNewClient(c.UCenterRPC)),
		AssetRPC:  uclient.NewAsset(zrpc.MustNewClient(c.UCenterRPC)),
		Kafka:     kcli,
	}
	s.Init()
	return s
}
