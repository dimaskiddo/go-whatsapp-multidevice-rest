package router

import (
	"time"

	cache "github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/memory"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
)

func HttpCacheInMemory(cap int, ttl int) *cache.Client {
	// Check if Cache Capacity is Zero or Less Than Zero
	if cap <= 0 {
		// Set Default Cache Capacity to 1000 Keys
		cap = 1000
	}

	// Create New InMemory Cache Adapter
	memcache, err := memory.NewAdapter(
		// Set InMemory Cache Adapter Algorithm to LRU
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(cap),
	)

	if err != nil {
		log.Print(nil).Error(err.Error())
		return nil
	}

	// Create New Cache Client with InMemory Adapter
	client, err := cache.NewClient(
		cache.ClientWithAdapter(memcache),
		cache.ClientWithTTL(time.Duration(ttl)*time.Second),
	)

	if err != nil {
		log.Print(nil).Error(err.Error())
		return nil
	}

	// Return Cache Client
	return client
}
