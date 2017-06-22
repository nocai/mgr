package roleser

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"mgr/models"
	"github.com/astaxie/beego/orm"
	"fmt"
)

func TestInsertRole(t *testing.T) {
	role, err := GetRoleById(1)
	fmt.Println(err)
	if err != nil {
		if err == orm.ErrNoRows {
			role = &models.Role{
				Id:1,
				RoleName:"roleName",
			}
			err = InsertRole(role)
			Convey("InsertRole", t, func() {
				Convey("err is nil", func() {
					So(err, ShouldBeNil)
				})
			})
		}
	}
	Convey("GetRoleById", t, func() {
		Convey("err is nil", func() {
			So(err, ShouldBeNil)
		})
	})

	exist, err := isExistOfRole(&models.Role{
		Id:role.Id,
		RoleName:role.RoleName,
	})
	Convey("isExistOfRole", t, func() {
		Convey("err is nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("not exist", func() {
			So(exist, ShouldBeFalse)
		})
	})


}


