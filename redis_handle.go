package cache

import "time"
import "github.com/gomodule/redigo/redis"

type RedisHandle struct {
	pool *redis.Pool
}

var (
	network string = "tcp"
	address string = "127.0.0.1:6379"
	password string = ""
	maxConnection int = 100

	channel string = "LocalStorageSync"
)

func (r *RedisHandle) Init() {
	r.pool := &redis.Pool{
		MaxIdle: maxConnection,
		MaxActive: maxConnection,
		Wait: false,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, error := c.Do("PING")
			return error
		},
		Dail: func() (redis.Conn, error) {
			return dial(network, address, password)
		},
	}

	return r
}

func dial() (redis.Conn, error) {
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