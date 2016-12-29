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
	beego.Debug(ctr.Input())

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)

	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	adminName := ctr.GetString("admin_name")

	data := make(map[string]interface{})
	data["admin_name"] = adminName

	key := util.NewPagerKey(page, rows, data, sort, order)
	pager := models.PageAdmin(key)
	ctr.Print(pager.Pagination)
}

func (ctr *AdminController) Post() {
	beego.Debug(ctr.Input())

	adminName := ctr.GetString("admin_name")
	password := ctr.GetString("password")

	admin := &models.Admin{AdminName:adminName, User:models.User{Username:adminName, Password:password}}
	err := models.AddAdmin(admin)
	if err != nil {
		ctr.PrintErrorMsg(err.Error())
		return
	}

	ctr.PrintOk()
}

func (ctr *AdminController) Delete() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug(fmt.Sprintf("id = %v", id))

	err := models.DeleteAdminById(id)
	if err != nil {
		ctr.PrintError()
		return
	}
	ctr.PrintOk()
}