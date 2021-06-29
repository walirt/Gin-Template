package libs

import (
	"context"
	"north-api/services"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var Cache *cache.Cache

func InitCache() {
	Cache = cache.New(&cache.Options{
		Redis: services.Rdb,
	})
}

func CacheSet(key string, value interface{}, ttl time.Duration, ctxs ...context.Context) error {
	ctx := context.TODO()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func CacheGet(key string, value interface{}, ctxs ...context.Context) error {
	ctx := context.TODO()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return Cache.Get(ctx, key, value)
}

func CacheScan(match string, cursor uint64, count int64, ctxs ...context.Context) *redis.ScanIterator {
	ctx := context.TODO()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return services.Rdb.Scan(ctx, cursor, match, 0).Iterator()
}

func CacheDelete(key string, ctxs ...context.Context) error {
	ctx := context.TODO()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return Cache.Delete(ctx, key)
}

func CacheExists(key string, ctxs ...context.Context) bool {
	ctx := context.TODO()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return Cache.Exists(ctx, key)
}
