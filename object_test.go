package aredis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	origin = "origin"
	object = "object"
)

func TestObject(t *testing.T) {
	Convey("Object", t, func() {
		c, err := New(getRedisURL(), NewDefaultConfig(name, version))
		So(err, ShouldBeNil)

		Convey("It shouldn't overwrite if object doesn't exist", func() {
			o := struct{ Name string }{"Name"}
			err := c.GetObject(origin, object, &o)
			So(err, ShouldBeNil)
			So(o.Name, ShouldEqual, "Name")
		})

		Convey("It should save and get object", func() {
			defer resetDb()

			err := c.SaveObject(origin, object, struct{ Name string }{"Name"})
			So(err, ShouldBeNil)

			n := struct{ Name string }{}
			err = c.GetObject(origin, object, &n)
			So(err, ShouldBeNil)
			So(n.Name, ShouldEqual, "Name")
		})

		Convey("It should save and get settings", func() {
			defer resetDb()

			err := c.SaveSettings(origin, struct{ Settings bool }{true})
			So(err, ShouldBeNil)

			s := struct{ Settings bool }{}
			err = c.GetSettings(origin, &s)
			So(err, ShouldBeNil)
			So(s.Settings, ShouldBeTrue)
		})
	})
}
