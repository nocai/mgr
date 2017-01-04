package controllers

import (
	"mgr/models"
	"github.com/astaxie/beego"
	"mgr/util"
)

type ResSelectController struct {
	BaseController
}

func (ctr *ResSelectController) Get() {
	var pid int64 = 0;
	resSelects, err := models.FindResByPid(pid)
	if err != nil {
		beego.Error(err)
	}

	var cb []util.Combobox
	for _, res := range resSelects {
		cb = append(cb, util.Combobox{Id:res.Id, Text:res.ResName})
	}
	ctr.Print(cb)
}