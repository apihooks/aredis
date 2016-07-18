package aredis

import (
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	seperatorKey = ":"
	settingsKey  = "settings"
)

// Config is the user definable options.
type Config struct {
	// Name is identifier of app using this library. This is used as prefix in
	// keys.
	Name string

	// Version is version of app using this library. This is used as prefix in
	// keys.
	Version string

	// MaxIdle is redigo/redis setting.
	MaxIdle int

	// MaxActive is redigo/redis setting.
	MaxActive int

	// IdleTimeout is redigo/redis setting.
	IdleTimeout time.Duration
}

// NewDefaultConfig returns *Config with sane defaults for redigo/redis.
func NewDefaultConfig(name, version string) *Config {
	return &Config{
		Name:        name,
		Version:     version,
		MaxIdle:     3,
		MaxActive:   0, // unlimited
		IdleTimeout: 240 * time.Second,
	}
}

// Client is a client for this library. Use New() to initalize it.
type Client struct {
	// Name is the identifier of worker using this library. This is used to
	// prefix keys along with Version in Redis.
	Name    string
	Version string

	// Seperator is a string that seperates Name, Version and the key.
	Seperator string

	// pool of redigo Redis connections.
	pool *redis.Pool
}

// New initializeds new Client along with Redis connection pool.
func New(url string, c *Config) (*Client, error) {
	pool := &redis.Pool{
		MaxIdle:     c.MaxIdle, // no. of idle conns in pool
		IdleTimeout: c.IdleTimeout,
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
		Name:      c.Name,
		Version:   c.Version,
		Seperator: seperatorKey,
		pool:      pool,
	}

	return client, nil
}

// Do is the primary method for interacting with Redis. It runs the given cmd
// in a redis conn from the pool and closes it when done.
//
// It also prefixes all keys with name and version passed in Config when
// initializing.
func (c *Client) Do(cmd, key string, rest ...interface{}) (interface{}, error) {
	conn := c.GetConn()
	defer conn.Close()

	prefix := []interface{}{c.Prefix(key)}
	prefix = append(prefix, rest...)

	return conn.Do(cmd, prefix...)
}

// WithOrigin prefixes origin to key to indicate key belongs to that origin.
func (c *Client) WithOrigin(origin, key string) string {
	return strings.Join([]string{origin, key}, seperatorKey)
}

// Prefix prefixes given key with name and version passed in Config when
// initializing.
func (c *Client) Prefix(key string) string {
	return strings.Join([]string{c.Name, c.Version, key}, seperatorKey)
}

// GetPool returns the redigo Redis connection pool.
func (c *Client) GetPool() *redis.Pool {
	return c.pool
}

// GetConn returns a single redigo Redis connection. Be sure to close the
// connection when finished.
func (c *Client) GetConn() redis.Conn {
	return c.pool.Get()
}

// Close closes Redis connection.
func (c *Client) Close() error {
	return c.pool.Close()
}

// IsErrNil returns true if err is key doesn't exist in Redis.
func (c *Client) IsErrNil(err error) bool {
	return err != nil && err == redis.ErrNil
}
