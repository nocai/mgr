package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"mgr/models/service/userser"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) T() {
	beego.Info("tttttttttttt")
	user := userser.User{Username: "username", Password: "password"}
	c.Data["json"] = user
	c.ServeJSON()
}

type HtmlController struct {
	BaseController
}

func (this *HtmlController) Get() {
	beego.Info(fmt.Sprintf("%#v", this.Ctx.Input.Params()))
	this.TplName = this.GetString(":splat")
}
