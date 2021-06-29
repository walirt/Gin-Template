package services

import "github.com/go-redis/redis/v8"

var Rdb *redis.Client

func OpenRedis(options *redis.Options) {
	Rdb = redis.NewClient(options)
}
