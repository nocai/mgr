package models

import (
	"testing"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"mgr/tests/testutil"
)

func TestUserKey(t *testing.T) {
	userKey := &UserKey{
		Key : testutil.GetTestKey(),
		User:&User{
			Id:1,
			Username:"username",
		},
	}
	sqler := userKey.NewSqler()
	Convey("TestUserKey", t, func() {
		countSql := sqler.GetCountSql()
		fmt.Println("countSql = ", countSql)
		Convey("the countSql should not be blank and should be contains \"select *\"", func() {
			So(countSql, ShouldNotBeBlank)
			So(countSql, ShouldContainSubstring, "select *")
		})
		sql := sqler.GetSql()
		fmt.Println("sql = ", sql)
		Convey("the sql should not be contains \"username\"", func() {
			So(sql, ShouldContainSubstring, "username")
		})
		args := sqler.GetArgs()
		fmt.Println(args)
		Convey("the time should be now", func() {
			So(args, ShouldNotBeEmpty)
			So(args, ShouldContain, "username")
		})
	})
}


