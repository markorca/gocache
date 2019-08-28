package cache

import "time"
import "github.com/gomodule/redigo/redis"

type RedisHandle struct {
	pool *redis.Pool
}

func (r *RedisHandle) Init() {
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
		Dail: func() (redis.Conn, error) {
			return dial(network, address, password)
		},
	}

	return r
}

func dial() (redis.Conn, error) {
	var network string = "tcp"
	var address string = "127.0.0.1:6379"
	var password string = ""

	c, error := redis.Dial(network, address)
	if err != nil {
		return nil, error
	}

	if password != "" {
		if _, error := c.Do("AUTH", password); error != nil {
			c.Close()
			return nil, error
		}
	}

	return c, error
}

func (r *RedisHandle) Publish(message string) error {
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