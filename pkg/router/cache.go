package router

import (
	"time"

	"github.com/labstack/echo/v4"

	cache "github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/memory"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
)

func HttpCacheInMemory(cap int, ttl int) echo.MiddlewareFunc {
	// Check if Cache Capacity is Zero or Less
	if cap <= 0 {
		// Set Default Cache Capacity
		cap = 1000
	}

	// Check if Cache TTL is Zero or Less
	if ttl <= 0 {
		// Set Default Cache TTL
		ttl = 5
	}

	// Create New In-Memory Cache Adapter
	memcache, err := memory.NewAdapter(
		// Set In-Memory Cache Adapter Algorithm to LRU
		// and With Desired Capacity
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(cap),
	)

	if err != nil {
		log.Print(nil).Error(err.Error())
		return nil
	}

	// Create New Cache
	cache, err := cache.NewClient(
		// Set Cache Adapter with In-Memory Cache Adapter
		// and Set Cache TTL in Second(s)
		cache.ClientWithAdapter(memcache),
		cache.ClientWithTTL(time.Duration(ttl)*time.Second),
	)

	if err != nil {
		log.Print(nil).Error(err.Error())
		return nil
	}

	// Return Cache as Echo Middleware
	return cache.Middleware()
}
