package controllers

import (
	"github.com/astaxie/beego"
	"mgr/models"
	"fmt"
)

type LoginController struct {
	BaseController
}

//func (ctr *LoginController) Get() {
//	ctr.TplName = "login.html"
//}

// 登陆
func (ctr *LoginController) Post() {
	username := ctr.GetString("username")
	password := ctr.GetString("password")

	admin, err := models.Login(username, password)
	if err != nil {
		beego.Debug(err.Error())
		ctr.PrintError(err)
		return
	}
	fmt.Println(admin)

	//var role *models.Role
	//role, err = models.GetRoleById(admin.RoleId)
	//if err != nil {
	//	beego.Error(err.Error())
	//	ctr.PrintError()
	//}
	//
	//beego.Debug(role)

	ctr.PrintOk();

}


