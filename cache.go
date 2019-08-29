package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	// handle memcache
	mem *MemcacheHandle

	// local storage interface
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

func (c *Cache) GetObject(key string) (interface{}, error) {
	cacheItem, err := c.lsIntf.Get(key)
	if err != nil {
		return nil, err
	}
	
	var value interface{}
	if err := json.Unmarshal(cacheItem.Value, &value); err != nil {
		return nil, err
	}

	return value, err
}

// need distributed
func (c *Cache) SetObject(key string, value interface{}, expiration int32) error {
	jsValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cacheItem := &CacheItem{
		Key: key,
		Value: jsValue,
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