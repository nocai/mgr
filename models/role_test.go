package models

import (
	"testing"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"mgr/tests/testutil"
)

func TestRoleKey(t *testing.T) {
	roleKey := &RoleKey{
		Key:testutil.GetTestKey(),
		Role:&Role{
			Id:1,
			RoleName:"roleName",
		},
	}
	sqler := roleKey.NewSqler()
	Convey("TestAdminKey", t, func() {
		countSql := sqler.GetCountSql()
		fmt.Println("countSql = ", countSql)
		Convey("the countSql should not be blank and should be contains \"select *\"", func() {
			So(countSql, ShouldNotBeBlank)
			So(countSql, ShouldContainSubstring, "select *")
		})
		sql := sqler.GetSql()
		fmt.Println("sql = ", sql)
		Convey("the sql should not be contains \"role_name\"", func() {
			So(sql, ShouldContainSubstring, "role_name")
		})
		args := sqler.GetArgs()
		fmt.Println(args)
		Convey("the time should be now", func() {
			So(args, ShouldNotBeEmpty)
		})
	})
}


