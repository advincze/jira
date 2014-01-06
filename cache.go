package jira

import (
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
	// log.Printf("[%s] caching put", key)

	//ioutil.WriteFile(time.Now().String(), []byte(fmt.Sprintf("%s", obj)), 0777)

	// log.Printf("[%s] caching acquiring write lock", key)
	c.Lock()
	defer c.Unlock()

	c.cacheMap[key] = &cachedObj{
		expire: time.Now().Add(dur),
		obj:    obj,
	}

	// log.Printf("[%s] caching releasing write lock", key)
}

func (c *cache) get(key string) (*cachedObj, bool) {
	// log.Printf("[%s] caching get", key)
	// log.Printf("[%s] caching acquiring read lock", key)
	c.RLock()
	defer c.RUnlock()

	cachedObj, ok := c.cacheMap[key]
	// log.Printf("[%s] caching releasing read lock", key)
	return cachedObj, ok
}

func (c *cache) getOrExpire(key string) interface{} {
	// log.Printf("[%s] caching getOrExpire", key)
	now := time.Now()
	cachedObj, ok := c.get(key)
	if ok {
		// log.Printf("[%s] found in cache", key)
		if cachedObj.expire.After(now) {
			// log.Printf("[%s] cache returning ", key)
			return cachedObj.obj
		} else {
			// log.Printf("[%s] cache expired", key)
			c.del(key)
		}
	}
	return nil
}

func (c *cache) del(key string) {
	// log.Printf("[%s] caching del", key)
	// log.Printf("[%s] caching acquiring write(del) lock", key)
	c.Lock()
	defer c.Unlock()

	delete(c.cacheMap, key)
	// log.Printf("[%s] caching releasing write(del) lock", key)
}

func (c *cache) getOrRunCacheAndReturn(dur time.Duration, key string, retrieveFn func() interface{}) interface{} {
	// log.Printf("[%s] caching getOrRunCacheAndReturn", key)
	if obj := c.getOrExpire(key); obj != nil {
		// log.Printf("[%s] getOrRunCacheAndReturn returning cached instance", key)
		return obj
	}
	// log.Printf("[%s] getOrRunCacheAndReturn retrieving instance from jira", key)
	obj := retrieveFn()
	// log.Printf("[%s] getOrRunCacheAndReturn caching retrieved instance", key)
	c.put(dur, key, obj)
	// log.Printf("[%s] getOrRunCacheAndReturn caching done returning", key)
	return obj
}
