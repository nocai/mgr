package controllers

import (
	"github.com/astaxie/beego"
	"mgr/conf"
	"mgr/util"
)

type ResController struct {
	BaseController
}

func (ctr *ResController) Get() {
	beego.Debug(ctr.Input())

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)

	resName := ctr.GetString("resName")
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	data := map[string]interface{}{
		"resName" : resName,
	}

	util.NewPagerKey(page, rows, data, sort, order)



}
