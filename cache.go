package jira

import (
	"log"
	"sync"
	"time"
)

type cache struct {
	sync.RWMutex
	cacheMap map[string]*cachedObj
}

type cachedObj struct {
	expire time.Time
	obj    interface{}
}

var defaultCache *cache

func init() {
	defaultCache = &cache{
		cacheMap: make(map[string]*cachedObj),
	}
}

func (c *cache) put(dur time.Duration, key string, obj interface{}) {
	log.Printf("[%s] caching", key)
	c.Lock()
	c.cacheMap[key] = &cachedObj{
		expire: time.Now().Add(dur),
		obj:    obj,
	}
	c.Unlock()
	log.Printf("[%s] caching finished", key)
}

func (c *cache) get(key string) interface{} {
	now := time.Now()
	log.Printf("looking for [%s] in the cache", key)
	c.RLock()
	if cachedObj, ok := c.cacheMap[key]; ok {
		c.RUnlock()
		log.Printf("found [%s] in the cache", key)
		if cachedObj.expire.After(now) {
			return cachedObj.obj
		} else {
			log.Printf("[%s] cache expired", key)
			c.Lock()
			delete(c.cacheMap, key)
			c.Unlock()
		}
	}
	c.RUnlock()
	return nil
}

func (c *cache) getOrRunCacheAndReturn(dur time.Duration, key string, retrieveFn func() interface{}) interface{} {
	if obj := c.get(key); obj != nil {
		log.Printf("[%s] returning cached instance", key)
		return obj
	}

	obj := retrieveFn()
	c.put(dur, key, obj)
	log.Printf("[%s] caching done returning", key)
	return obj
}
