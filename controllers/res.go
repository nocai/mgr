package controllers

import (
	"github.com/astaxie/beego"
	"mgr/conf"
	"mgr/util"
	"mgr/models"
	"fmt"
)

type ResController struct {
	BaseController
}

func (ctr *ResController) Get() {
	beego.Debug(ctr.Input())

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)

	resName := ctr.GetString("resName")
	path := ctr.GetString("path")
	pid, _ := ctr.GetInt64("pid", 0)

	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	data := map[string]interface{}{
		"resName" : resName,
		"path" : path,
		"pid" : pid,
	}
	key := util.NewPagerKey(page, rows, data, sort, order)
	pager, err := models.PageRes(key)
	if err != nil {
		beego.Error(err)
	}
	ctr.Print(pager.Pagination)
}

func (ctr *ResController) Post() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug(fmt.Sprintf("id = %+v", id))

	pid, _ := ctr.GetInt64("pid", 0)
	resName := ctr.GetString("res_name")
	path := ctr.GetString("path")

	res := models.Res{Id:id, ResName:resName, Path:path, Pid:pid}
	var err error
	if id == 0 {
		// 添加
		err = models.InsertRes(&res)
	} else {
		// 修改
		err = models.UpdateRes(&res)
	}

	if err != nil {
		ctr.PrintErrorMsg(err.Error())
		return
	}
	ctr.PrintOk()
}

func (ctr *ResController) Delete() {
	id, _ :=ctr.GetInt64(":id", 0)
	beego.Debug("id = %v", id)

	err := models.DeleteResById(id)
	if err != nil {
		ctr.PrintErrorMsg(err.Error())
		return
	}
	ctr.PrintOk()
}