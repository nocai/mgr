package models

import (
	"testing"
	"time"
	"math/rand"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"mgr/util/key"
	"mgr/conf"
)

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


