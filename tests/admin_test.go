package test

import (
	"testing"
	"mgr/models"
	"mgr/util"
)

func TestAdmin(t *testing.T) {
	adminName := string(util.Krand(6, util.KC_RAND_KIND_UPPER))
	password := string(util.Krand(6, util.KC_RAND_KIND_UPPER))

	admin := &models.Admin{AdminName:adminName, User:models.User{Username:adminName, Password:password}}
	t.Logf("%+v",admin)
	// 测试添加
	err := models.InsertAdmin(admin)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v",admin)

	// 再添加一次,用户名重复不
	err = models.InsertAdmin(admin)
	if err != models.ErrInsert {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v",admin)

	// 取数据
	admin, err = models.GetAdminById(admin.Id)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v",admin)

	admin, err = models.GetAdminByUserId(admin.UserId, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v",admin)

	newPassword := string(util.Krand(6, util.KC_RAND_KIND_UPPER))
	admin.User.Password = newPassword
	t.Logf("%+v",admin)

	admin.User.Id = admin.UserId
	err = models.UpdateAdmin(admin)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if newPassword != admin.User.Password {
		t.Error("修改失败")
		t.FailNow()
	}
	t.Logf("%+v",admin)

	err = models.DeleteAdminById(admin.Id)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v",admin)

	admin, err = models.GetAdminById(admin.Id)
	if err != models.ErrQuery {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v",admin)

}

