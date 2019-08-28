package cache

import "github.com/bradfitz/gomemcache/memcache"

type MemcacheHandle struct {
	client *memcache.Client
}

func NewMemcacheHandle() *MemcacheHandle {
	m := &MemcacheHandle{}
	m.Init()
	return m
}

func cacheItemToItem(cacheItem *CacheItem) *memcache.Item {
	return &memcache.Item{
		Key: cacheItem.Key,
		Value: cacheItem.Value,
		Expiration: cacheItem.Expiration,
	}
}

func itemToCacheItem(item *memcache.Item) *CacheItem {
	return &CacheItem{
		Key: item.Key,
		Value: item.Value,
		Expiration: item.Expiration,
	}
}

func (m *MemcacheHandle) Init() *MemcacheHandle {
	var address string = "127.0.0.1:11211"

	m.client = memcache.New(address)
	return m
}

func (m *MemcacheHandle) Set(cacheItem *CacheItem) error {
	item := cacheItemToItem(cacheItem)

	if err := m.client.Set(item); err != nil {
		return err
	}

	return nil
}

func (m *MemcacheHandle) Get(key string) (*CacheItem, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}

	cacheItem := itemToCacheItem(item)

	return cacheItem, nil
}