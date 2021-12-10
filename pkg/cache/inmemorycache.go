package cache

import "sync"

type inMemoryCache struct {
	// 保存键值对
	c map[string][]byte
	// 对map的并发访问提供读写锁，支持多个goroutine读写
	mutex sync.RWMutex
	// 记录缓存状态，内嵌结构体，组合的方式
	Stat
}

func (cache *inMemoryCache) Set(key string, value []byte) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	temp, exist := cache.c[key]
	if exist {
		cache.del(key, temp)
	}
	cache.c[key] = value
	//更新状态里面的信息
	cache.Stat.add(key, value)
	return nil
}

func (cache *inMemoryCache) Get(key string) ([]byte, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	return cache.c[key], nil
}

func (cache *inMemoryCache) Del(key string) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	v, exist := cache.c[key]
	if exist {
		delete(cache.c, key)
		cache.Stat.del(key, v)
	}
	return nil
}

func (cache *inMemoryCache) GetStat() Stat {
	return cache.Stat
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		c:     make(map[string][]byte),
		mutex: sync.RWMutex{},
		Stat:  Stat{},
	}
}
