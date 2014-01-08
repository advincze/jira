package jira

import (
	"sync"
	"time"
)

var defaultCache = &cache{cacheMap: make(map[string]*cachedObj)}

type cache struct {
	sync.RWMutex
	cacheMap map[string]*cachedObj
}

type cachedObj struct {
	expire time.Time
	obj    interface{}
}

func (c *cache) put(dur time.Duration, key string, obj interface{}) {
	c.Lock()
	defer c.Unlock()

	c.cacheMap[key] = &cachedObj{
		expire: time.Now().Add(dur),
		obj:    obj,
	}
}

func (c *cache) get(key string) (*cachedObj, bool) {
	c.RLock()
	defer c.RUnlock()

	cachedObj, ok := c.cacheMap[key]
	return cachedObj, ok
}

func (c *cache) getOrExpire(key string) interface{} {
	now := time.Now()
	cachedObj, ok := c.get(key)
	if ok {
		if cachedObj.expire.After(now) {
			return cachedObj.obj
		} else {
			c.del(key)
		}
	}
	return nil
}

func (c *cache) del(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.cacheMap, key)
}

func (c *cache) getOrRunCacheAndReturn(dur time.Duration, key string, retrieveFn func() interface{}) interface{} {
	if obj := c.getOrExpire(key); obj != nil {
		return obj
	}
	obj := retrieveFn()
	c.put(dur, key, obj)
	return obj
}
