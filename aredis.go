package aredis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var SettingsKey = "settings"

type Client struct {
	// Name is the identifier of worker using this library. This is used to
	// prefix keys along with Version in Redis.
	Name    string
	Version string

	pool *redis.Pool
}

// New initializeds new Client along with Redis connection pool.
func New(url, name, version string) (*Client, error) {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", url) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := pool.Get()
	defer conn.Close()

	// ping to make sure connection works
	if _, err := conn.Do("PING"); err != nil {
		return nil, err
	}

	return &Client{pool: pool}, nil
}
