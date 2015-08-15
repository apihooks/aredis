package aredis

import (
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	SettingsKey  = "settings"
	SeperatorKey = ":"
)

type Client struct {
	// Name is the identifier of worker using this library. This is used to
	// prefix keys along with Version in Redis.
	Name    string
	Version string

	// Seperator is a string that seperates Name, Version and the key.
	Seperator string

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

	client := &Client{
		Name: name, Version: version, Seperator: SeperatorKey, pool: pool,
	}

	return client, nil
}

// Do is the primary method for interacting with Redis. It prefixes all keys.
func (c *Client) Do(cmd, key string, rest ...interface{}) (interface{}, error) {
	conn := c.GetConn()
	defer conn.Close()

	prefix := []interface{}{c.Prefix(key)}
	prefix = append(prefix, rest...)

	return conn.Do(cmd, prefix...)
}

// Prefix adds Name and Version prefixes to key.
func (c *Client) Prefix(key string) string {
	return strings.Join([]string{c.Name, c.Version, key}, SeperatorKey)
}

// GetPool returns the redigo redis connection pool.
func (c *Client) GetPool() *redis.Pool {
	return c.pool
}

// GetConn returns a single redigo redis connection. It's upto the caller to
// close the connection when done with it.
func (c *Client) GetConn() redis.Conn {
	return c.pool.Get()
}

// Close closes redis connections.
func (c *Client) Close() error {
	return c.pool.Close()
}
