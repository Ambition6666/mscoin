package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/ucenter/client"
	"jobcenter/internal/config"
	"jobcenter/internal/database"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *database.MongoClient
	KafkaClient *database.KafkaClient
	RedisClient *database.RedisClient
	AssetRPC    client.Asset
}

func NewServiceContext(c config.Config) *ServiceContext {
	kcli := database.NewKafkaClient(c.Kafka)
	kcli.StartWrite()

	rcli := database.NewRedisClient(c.CacheRedis)
	return &ServiceContext{
		Config:      c,
		MongoClient: database.ConnectMongo(c.Mongo),
		KafkaClient: kcli,
		RedisClient: rcli,
		AssetRPC:    client.NewAsset(zrpc.MustNewClient(c.UCenterRPC)),
	}
}
