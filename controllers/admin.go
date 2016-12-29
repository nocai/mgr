package controllers

import (
	"github.com/astaxie/beego"
	"mgr/conf"
	"mgr/util"
	"mgr/models"
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
