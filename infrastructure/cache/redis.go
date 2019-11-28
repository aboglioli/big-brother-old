package cache

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/cache"
	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	client *redis.Client
}

func Redis() cache.Cache {
	conf := config.Get()
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisURL,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})
	return &redisCache{client}
}

func (r *redisCache) Get(k string) interface{} {
	v, err := r.client.Get(k).Result()
	if err != nil {
		return nil
	}
	return v
}

func (r *redisCache) Set(k string, v interface{}, d time.Duration) {
	r.client.Set(k, v, d)
}
