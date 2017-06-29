package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
)

type HtmlController struct {
	BaseController
}

func (this *HtmlController) Get() {
	beego.Info(fmt.Sprintf("%#v",this.Ctx.Input.Params()))
	this.TplName = this.GetString(":splat")
}

