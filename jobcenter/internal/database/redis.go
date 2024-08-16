package database

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RedisClient struct {
	Cache cache.Cache
}

func NewRedisClient(cc redis.RedisConf) *RedisClient {
	newRedis := redis.MustNewRedis(cc)
	node := cache.NewNode(newRedis, nil, cache.NewStat("jobcenter"), nil)
	return &RedisClient{
		Cache: node,
	}
}
