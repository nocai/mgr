package userser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUsernamPassMatched(t *testing.T) {
	matched, err := UsernamPassMatched("b", "b")
	Convey("TestUsernamPassMatched", t, func() {
		Convey("err is nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("matched", func() {
			So(matched, ShouldBeTrue)
		})
	})
}
