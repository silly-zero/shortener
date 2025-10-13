package cache

import (
	"sync"
	"time"
)

// LocalCache 本地缓存实现
type LocalCache struct {
	cache map[string]cacheItem
	lock  sync.RWMutex
	ttl   time.Duration
}

type cacheItem struct {
	value      string
	expireTime time.Time
}

// NewLocalCache 创建本地缓存实例
func NewLocalCache(ttl time.Duration) *LocalCache {
	c := &LocalCache{
		cache: make(map[string]cacheItem),
		ttl:   ttl,
	}

	// 启动清理过期项的后台协程
	go c.cleanupLoop()

	return c
}

// Get 获取缓存项
func (c *LocalCache) Get(key string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, found := c.cache[key]
	if !found {
		return "", false
	}

	// 检查是否过期
	if time.Now().After(item.expireTime) {
		return "", false
	}

	return item.value, true
}

// Set 设置缓存项
func (c *LocalCache) Set(key, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache[key] = cacheItem{
		value:      value,
		expireTime: time.Now().Add(c.ttl),
	}
}

// cleanupLoop 定期清理过期缓存项
func (c *LocalCache) cleanupLoop() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.lock.Lock()
			now := time.Now()
			for k, v := range c.cache {
				if now.After(v.expireTime) {
					delete(c.cache, k)
				}
			}
			c.lock.Unlock()
		}
	}
}