package svc

import (
	"common/msdb"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"market/database"
	"market/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	Cache       cache.Cache
	Db          *msdb.MsDB
	MongoClient *database.MongoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	newRedis := c.CacheRedis[0].NewRedis()
	db := database.ConnMysql(c.Mysql.DataSource)
	node := cache.NewNode(newRedis, nil, nil, nil)
	return &ServiceContext{
		Config:      c,
		Cache:       node,
		Db:          db,
		MongoClient: database.ConnectMongo(c.Mongo),
	}
}
