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

	// redis conn, use it to update/delete
	rds *RedisHandle
}

func Init() *Cache {
	c := &Cache{}

	c.mem = NewMemcacheHandle()
	c.rds = NewRedisHandle()

	c.lsIntf = LocalStorage(c.mem)

	go subscribe(c)

	return c
}

// use redis subscribe/publish function to sync object
func subscribe(c *Cache) error {
	var channel string = "LocalStorageSync"

	conn := c.rds.pool.Get()
	// defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(channel); err != nil {
		return err
	}

	go func() {
		for {
			switch v:= psc.Receive().(type) {
			case redis.Message:
				// v.Data, v.Channel
				if (v.Channel == channel) {
					if err := LocalStorageSync(c, v.Data); err != nil {
						panic(err)
					}
				}
			case redis.Subscription:
				// NOTHING
				// fmt.Println(v)
			case error:
				// TODO error
				panic(v)
			}
		}
	}()

	return nil
}

func LocalStorageSync(c *Cache, data []byte) error {
	var cacheItem *CacheItem
	if err := json.Unmarshal(data, &cacheItem); err != nil {
		return err
	}

	method := cacheItem.Method
	switch method {
	case "SET" :
		if err := c.lsIntf.Set(cacheItem); err != nil {
			return err
		}
	case "DELETE":
		if err := c.lsIntf.Delete(cacheItem.Key); err != nil {
			return err
		}
	}

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

func (c *Cache) SetObject(key string, value interface{}, expiration int32) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cacheItem := &CacheItem{
		Key: key,
		Value: jsonValue,
		Expiration: expiration,
		Method: "SET",
	}

	if err := c.lsIntf.Set(cacheItem); err != nil {
		return err
	}

	jsonCacheItem, _ := json.Marshal(cacheItem)
	if err := c.rds.Publish(jsonCacheItem); err != nil {
		return err
	}

	return nil
}

func (c *Cache) DeleteObject(key string) error {
	cacheItem := &CacheItem{
		Key: key,
		Method: "DELETE",
	}

	if err := c.lsIntf.Delete(key); err != nil {
		return err
	}

	jsonCacheItem, _ := json.Marshal(cacheItem)
	if err := c.rds.Publish(jsonCacheItem); err != nil {
		return err
	}

	return nil
}