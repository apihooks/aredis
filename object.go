package aredis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

// GetSettings gets settings for origin.
// key: <worker>:<version>:<origin>:settings
func (c *Client) GetSettings(origin string, s interface{}) error {
	return c.GetObject(origin, settingsKey, &s)
}

// SaveSettings saves given settings for origin.
// key: <worker>:<version>:<origin>:settings
func (c *Client) SaveSettings(origin string, s interface{}) error {
	return c.SaveObject(origin, settingsKey, &s)
}

// GetObject gets marshalled object for origin and key.
// key: <worker>:<version>:<origin>:<key>
func (c *Client) GetObject(origin, key string, o interface{}) error {
	raw, err := redis.Bytes(c.Do("GET", c.WithOrigin(origin, key)))
	if err != nil && !c.IsErrNil(err) {
		return err
	}

	if c.IsErrNil(err) {
		return nil
	}

	return json.Unmarshal(raw, &o)
}

// SaveObject saves object as json string for origin and key. It uses string
// instead of bytes for ease of debugging.
//
// key: <worker>:<version>:<origin>:<key>
func (c *Client) SaveObject(origin, key string, o interface{}) error {
	raw, err := json.Marshal(o)
	if err != nil {
		return err
	}

	_, err = c.Do("SET", c.WithOrigin(origin, key), string(raw))
	return err
}
