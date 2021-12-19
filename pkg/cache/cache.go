package cache

type Cache interface {
	// Set 新增缓存
	Set(string, []byte) error
	// Get 获取缓存
	Get(string) ([]byte, error)
	// Del 删除缓存
	Del(string) error
	// GetStat 获取缓存状态
	GetStat() Stat
}

type cacheType int

const (
	InMemory cacheType = iota
)

func New(typ string) Cache {
	var cache Cache
	if typ == "inMemoryCache" {
		cache = newInMemoryCache()
	}
	if typ == "rocksdb" {
		cache = newRocksdbCache()
	}
	if cache == nil {
		panic("unknown cache type : " + typ)
	}
	return cache
}
