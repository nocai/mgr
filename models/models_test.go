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


func TestAdminKey(t *testing.T) {
	now := time.Now()
	adminKey := &AdminKey{
		Key : key.New(conf.Page, conf.Rows, []string {}, []string{}, true),
		Admin:&Admin{
			ModelBase : ModelBase{
				CreateTime:time.Now(),
				UpdateTime:time.Now(),
			},
			Id:rand.Int63n(100),
			AdminName:"A",
			UserId:rand.Int63n(100),
			Invalid:ValidEnum_Invalid,
		},
		CreateTimeStart:now,
		UpdateTimeEnd:now,
		UpdateTimeStart:now,
		CreateTimeEnd:now,
		KeyWord:"keyword",
	}
	sqler := adminKey.NewSqler()
	Convey("TestAdminKey", t, func() {
		countSql := sqler.GetCountSql()
		fmt.Println("countSql = ", countSql)
		Convey("the countSql should not be blank and should be contains \"select *\"", func() {
			So(countSql, ShouldNotBeBlank)
			So("select *", ShouldNotBeIn, countSql)
		})
		sql := sqler.GetSql()
		fmt.Println("sql = ", sql)
		Convey("the sql should not be blank", func() {
			So(sql, ShouldNotBeBlank)
		})
		args := sqler.GetArgs()
		fmt.Println(args)
		Convey("the time should be now", func() {
			So(now, ShouldBeIn, args)
		})
	})
}

