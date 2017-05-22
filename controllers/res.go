package controllers
//
//import (
//	"github.com/astaxie/beego"
//	"mgr/conf"
//	"mgr/util"
//	"mgr/models"
//	"fmt"
//)
//
//type ResController struct {
//	BaseController
//}
//
//func (ctr *ResController) Get() {
//	beego.Debug(ctr.Input())
//
//	page, _ := ctr.GetInt64("page", conf.Page)
//	rows, _ := ctr.GetInt64("rows", conf.Rows)
//
//	resName := ctr.GetString("resName")
//	path := ctr.GetString("path")
//	pid, _ := ctr.GetInt64("pid", 0)
//	seq, _ := ctr.GetInt("seq", 0)
//
//	sort := ctr.GetString("sort", "id")
//	order := ctr.GetString("order", "desc")
//
//	key := &models.ResKey{}
//	key.Key = util.NewKey(page, rows, []string{sort}, []string{order}, true)
//	key.Res = models.Res{ResName:resName, Path:path, Pid:pid, Seq:seq}
//	pager := models.PageRes(key)
//	ctr.Print(pager)
//}
//
//func updateRes(id, pid int64, resName, path string) error {
//	res, err := models.GetResByResId(id)
//	if err != nil {
//		beego.Error(err)
//		return err
//	}
//	res.Pid = pid
//	res.ResName = resName
//	res.Path = path
//	return models.UpdateRes(res)
//}
//
//func (ctr *ResController) Post() {
//	beego.Debug(ctr.Input())
//
//	id, _ := ctr.GetInt64(":id", 0)
//	beego.Debug(fmt.Sprintf("id = %+v", id))
//
//	pid, _ := ctr.GetInt64("pid", models.Pid_Default)
//	resName := ctr.GetString("res_name")
//	path := ctr.GetString("path")
//
//	if id == 0 {
//		// 添加
//		res := models.Res{Id:id, ResName:resName, Path:path, Pid:pid}
//		ctr.PrintError(models.InsertRes(&models.ResVo{Res:res}))
//	} else {
//		// 修改
//		ctr.PrintError(updateRes(id, pid, resName, path))
//	}
//}
//
//func (ctr *ResController) Delete() {
//	id, _ := ctr.GetInt64(":id", 0)
//	beego.Debug("id = %v", id)
//
//	err := models.DeleteResById(id)
//	ctr.PrintError(err)
//}
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