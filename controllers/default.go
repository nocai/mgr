package controllers

import (
	"github.com/astaxie/beego"
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
	c.PrintOkMsgData("操作成功", "this is t method")

}
