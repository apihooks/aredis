package aredis

import (
	"net"
	"testing"

	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	redisUrl  = "localhost:6379"
	redisConn redis.Conn

	name    = "aredis"
	version = "0.1"
)

func init() {
	var err error

	if redisConn, err = redis.Dial("tcp", redisUrl); err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	Convey("It should return error if redis isn't reachable", t, func() {
		// start a server on port 0, so OS'll assign an random empty port
		l, _ := net.Listen("tcp", ":0")
		l.Close()

		_, err := New(l.Addr().String(), NewDefaultConfig(name, version))
		So(err, ShouldNotBeNil)
	})
}

func resetDb() {
	deleteScript := "return redis.call('del', unpack(redis.call('keys', ARGV[1])))"
	script := redis.NewScript(0, deleteScript)

	script.Do(redisConn, name+"*")
}
