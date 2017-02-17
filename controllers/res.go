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

	pid, _ := ctr.GetInt64("pid", -1)
	resName := ctr.GetString("res_name")
	path := ctr.GetString("path")

	var err error
	if id == 0 {
		// 添加
		res := models.Res{Id:id, ResName:resName, Path:path, Pid:pid}
		err = models.InsertRes(&models.ResVo{Res:res})
	} else {
		// 修改
		key := &models.ResKey{Res:models.Res{Id:id}}
		var ress []models.Res
		ress, err = models.FindResByKey(key)
		if err == nil && len(ress) > 0 {
			res := ress[0]
			res.ResName = resName
			res.Path = path
			res.Pid = pid
			err = models.UpdateRes(&res)
		}
	}

	if err != nil {
		ctr.PrintErrorMsg(err.Error())
		return
	}
	ctr.PrintOk()
}

func (ctr *ResController) Delete() {
	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = %v", id)

	err := models.DeleteResById(id)
	if err != nil {
		ctr.PrintErrorMsg(err.Error())
		return
	}
	ctr.PrintOk()
}

type ResSelectController struct {
	BaseController
}

func (ctr *ResSelectController) Get() {
	key := &models.ResKey{Res:models.Res{Pid:-1}}
	resSelects, err := models.FindResByKey(key)
	if err != nil {
		beego.Error(err)
	}

	var cb []util.Combobox
	for _, res := range resSelects {
		cb = append(cb, util.Combobox{Id:res.Id, Text:res.ResName})
	}
	ctr.Print(cb)
}