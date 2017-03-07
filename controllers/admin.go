package controllers

import (
	"github.com/astaxie/beego"
	"mgr/conf"
	"mgr/util"
	"mgr/models"
	"fmt"
)

type AdminController struct {
	BaseController
}

func (ctr *AdminController) Get() {
	ctr.debugInput()

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)

	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	adminName := ctr.GetString("admin_name")
	key := util.NewKey(page, rows, []string{sort}, []string{order}, true)
	pager, err := models.PageAdmin(&models.AdminKey{Key:key, Admin:models.Admin{AdminName:adminName}})
	if err != nil {
		beego.Error(err)
	}
	ctr.Print(pager.Pagination)
}

func addAdmin(adminName, username, password string) error {
	user := models.User{Username:username, Password:password}
	admin := models.Admin{AdminName:adminName}
	adminVo := &models.AdminVo{Admin:admin, User:user}
	return models.InsertAdminVo(adminVo)
}

func updateAdmin(id int64, adminName, username, password string) error {
	admin, err := models.GetAdminById(id)
	if err != nil {
		beego.Error(err)
		return err
	}
	admin.AdminName = adminName
	admin.User.Username = username
	admin.User.Password = password
	return models.UpdateAdmin(admin)
}

// 添加 修改
func (ctr *AdminController) Post() {
	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug(fmt.Sprintf("id = %v", id))
	beego.Debug(ctr.Input())

	adminName := ctr.GetString("admin_name")
	username := ctr.GetString("username")
	password := ctr.GetString("password")

	if id == 0 {// 添加
		ctr.PrintError(addAdmin(adminName, username, password))
	} else {// 更新
		ctr.PrintError(updateAdmin(id, adminName, username, password))
	}
}

func (ctr *AdminController) Delete() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug(fmt.Sprintf("id = %v", id))

	err := models.DeleteAdminById(id)
	ctr.PrintError(err)
}