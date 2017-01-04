package controllers

import (
	"mgr/models"
	"github.com/astaxie/beego"
)

type ResSelectController struct {
	BaseController
}

func (ctr *ResSelectController) Get() {
	var pid int64 = 0;
	resSelects, err := models.FindResSelectByPid(pid)
	if err != nil {
		beego.Error(err)
	}
	ctr.Print(resSelects)
}