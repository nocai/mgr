package controllers

import (
	"github.com/astaxie/beego"
	"mgr/conf"
	"mgr/models/service/resser"
	"mgr/util/key"
)

type ResController struct {
	BaseController
}

func (ctr *ResController) Get() {
	defer ctr.recoverPanic()

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)

	resName := ctr.GetString("resName")
	path := ctr.GetString("path")
	pid, _ := ctr.GetInt64("pid", 0)
	seq, _ := ctr.GetInt("seq", 0)

	sort := ctr.GetString("sort", "id")
	order := ctr.GetString("order", "desc")

	key := &resser.ResKey{
		Key: key.New(page, rows, []string{sort}, []string{order}),
		Res: resser.Res{
			ResName: resName,
			Path:    path,
			Pid:     pid,
			Seq:     seq,
		},
	}
	pager := resser.PageRes(key)
	ctr.PrintJson(pager)
}

func (ctr *ResController) Post() {
	defer ctr.recoverPanic()

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = ", id)

	pid, _ := ctr.GetInt64("pid", resser.Pid_Default)
	resName := ctr.GetString("res_name")
	path := ctr.GetString("path")

	if id == 0 {
		// 添加
		res := &resser.Res{Id: id, ResName: resName, Path: path, Pid: pid}
		ctr.PrintJson(resser.InsertRes(res))
	} else {
		// 修改
		res := resser.GetResByResId(id)
		res.Pid = pid
		res.ResName = resName
		res.Path = path
		ctr.PrintJson(resser.UpdateRes(res))
	}
}

func (ctr *ResController) Delete() {
	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = %v", id)

	err := resser.DeleteResById(id)
	ctr.PrintJson(err)
}

//
//type ResSelectController struct {
//	BaseController
//}
//
//func (ctr *ResSelectController) Get() {
//	key := &models.ResKey{Res:models.Res{Pid:-1}}
//	resSelects, err := models.FindResByKey(key)
//	if err != nil {
//		beego.Error(err)
//	}
//
//	var cb []util.Combobox
//	for _, res := range resSelects {
//		cb = append(cb, util.Combobox{Id:res.Id, Text:res.ResName})
//	}
//	ctr.Print(cb)
//}
