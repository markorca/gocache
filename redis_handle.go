package cache

import "time"
import "github.com/gomodule/redigo/redis"

type RedisHandle struct {
	pool *redis.Pool
}

func NewRedisHandle() *RedisHandle {
	r := &RedisHandle{}
	r.Init()
	return r
}

func (r *RedisHandle) Init() *RedisHandle {
	var network string = "tcp"
	var address string = "127.0.0.1:6379"
	var password string = ""
	var maxConnection int = 100

	r.pool = &redis.Pool{
		MaxIdle: maxConnection,
		MaxActive: maxConnection,
		Wait: false,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial(network, address, password)
		},
	}

	return r
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}

	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}

	return c, err
}

func (r *RedisHandle) Publish(message []byte) error {
	var channel string = "LocalStorageSync"

	conn := r.pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	if _, err := conn.Do("PUBLISH", channel, message); err != nil {
		return err
	}

	return nil
}