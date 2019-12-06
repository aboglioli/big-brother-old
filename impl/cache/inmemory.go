package cache

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/cache"
	gocache "github.com/patrickmn/go-cache"
)

const (
	NoExpiration      = gocache.NoExpiration
	DefaultExpiration = gocache.DefaultExpiration
)

type goCache struct {
	cache *gocache.Cache
}

func InMemory() cache.Cache {
	c := gocache.New(2*time.Minute, 5*time.Minute)

	return &goCache{c}
}

func (c *goCache) Get(k string) interface{} {
	data, ok := c.cache.Get(k)
	if !ok {
		return nil
	}

	return data
}

func (c *goCache) Set(k string, v interface{}, d time.Duration) {
	c.cache.Set(k, v, d)
}

func (c *goCache) Delete(k string) {
	c.cache.Delete(k)
}
