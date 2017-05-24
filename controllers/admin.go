package controllers

import (
	"mgr/conf"
	"mgr/models"
	"mgr/models/service/admin"
	"mgr/util/key"

	"github.com/astaxie/beego"
)

type AdminController struct {
	BaseController
}

//
func (ctr *AdminController) Get() {
	ctr.debugInput()

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")
	adminName := ctr.GetString("admin_name")

	key := key.New(page, rows, []string{sort}, []string{order}, true)
	_admin := &models.Admin{AdminName: adminName}
	pager, err := admin.PageAdmin(&models.AdminKey{Key: key, Admin: _admin})
	ctr.PrintResult(pager, err)
}

//func updateAdmin(id int64, adminName, username, password string) error {
//	admin, err := models.GetAdminById(id)
//	if err != nil {
//		beego.Error(err)
//		return err
//	}
//	admin.AdminName = adminName
//	admin.User.Username = username
//	admin.User.Password = password
//	return models.UpdateAdmin(admin)
//}

// 添加 修改
func (ctr *AdminController) Post() {
	ctr.debugInput()

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = ", id)

	adminName := ctr.GetString("admin_name")
	password := ctr.GetString("password")

	if id == 0 {
		// 添加
		adminVo := &models.AdminVo{Admin:  &models.Admin{AdminName: adminName}, User: &models.User{Username: adminName, Password: password}}
		ctr.PrintError(admin.InsertAdminVo(adminVo))
	} else {
		// 更新
		//ctr.PrintError(updateAdmin(id, adminName, username, password))
	}
}

//func (ctr *AdminController) Delete() {
//	beego.Debug(ctr.Input())
//
//	id, _ := ctr.GetInt64(":id", 0)
//	beego.Debug(fmt.Sprintf("id = %v", id))
//
//	err := models.DeleteAdminById(id)
//	ctr.PrintError(err)
//}
