package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	// handle memcache
	mem *MemcacheHandle

	lsIntf LocalStorage

	// redis conn, use it to update/delete in multi
	rds *RedisHandle
}

func Init() *Cache {
	c := &Cache{}

	c.mem = NewMemcacheHandle()
	c.rds = NewRedisHandle()

	c.lsIntf = LocalStorage(c.mem)

	go subscribe(c.rds)

	return c
}

// use redis subscribe/publish function to sync object
func subscribe(rds *RedisHandle) error {
	var channel string = "LocalStorageSync"

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
				fmt.Println(v.Data)
				// TODO v.Data, v.Channel
			case redis.Subscription:
				// NOTHING
			case error:
				// TODO error
			}
		}
	}()

	return nil
}

func (c *Cache) GetObject(key string) (*CacheItem, error) {
	cacheItem, err := c.lsIntf.Get(key)
	if err != nil {
		return nil, err
	}
	return cacheItem, err
}

// need distributed
func (c *Cache) SetObject(key string, value []byte, expiration int32) error {
	cacheItem := &CacheItem{
		Key: key,
		Value: value,
		Expiration: expiration,
	}

	if err := c.lsIntf.Set(cacheItem); err != nil {
		return err
	}

	return nil
}

// need distributed
func (c *Cache) DeleteObject() {

}