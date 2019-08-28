package cache

type Cache struct {
	// handle memcache
	mem *MemcacheHandle

	lsIntf *LocalStorage

	// redis conn, use it to update/delete in multi
	rds *RedisHandle
}

var (
	channel string = "LocalStorageSync"
)

func Init() {
	c := &Cache{}

	c.mem := MemcacheHandle.Init()
	c.rds := RedisHandle.Init()

	c.lsIntf := LocalStorage(c.mem)

	subscribe(c.rds)

	return c
}

// use redis subscribe/publish function to sync object
func subscribe(rds *RedisHandle) {
	conn := rds.pool.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(channel); err != nil {
		return err
	}

	go func() {
		for {
			switch v:= psc.Receive().(type) {
			case redis.Message:
				// TODO v.Data, v.Channel
			case redis.Subscription:
				// NOTHING
			case error:
				return v
			}
		}
	}
}

func (c *Cache) GetObject(key string) *CacheItem {
	cacheItem := c.lsIntf.Get(key)
	return cacheItem
}

// need distributed
func (c *Cache) SetObject(key string, value []byte, expiration int) {
	cacheItem := &cacheItem{
		Key: key,
		Value: value,
		Expiration: expiration,
	}

	c.lsIntf.Set(cacheItem)
}

// need distributed
func (c *Cache) DeleteObject() {

}