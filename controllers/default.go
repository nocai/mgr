package controllers

import (
	"github.com/astaxie/beego"
	"mgr/models"
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
	user:= models.User{Username:"username",Password:"password"}
	c.Data["json"] = user
	c.ServeJSON()
	//c.PrintOkMsgData("操作成功", "this is t method")

}
