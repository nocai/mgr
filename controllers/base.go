package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
)

const (
	OK int = 200
	Bad int = 400

	// session key "adminId"
	ADMIN_ID = "adminId"
	// session key "adminName"
	ADMIN_NAME = "adminName"
	// session key "auths"
	AUTHS = "auths"
)

type BaseController struct {
	beego.Controller
}

func (controller *BaseController) debugInput() {
	beego.Debug("输入参数:", controller.Input())
}
// 将用户信息放到Session里
//func (controller *BaseController) SetSesseionAdmin(admin models.Admin) {
//	controller.SetSession(ADMIN_ID, admin.Id)
//	controller.SetSession(ADMIN_NAME, admin.Username)
//}
//
//// 将用户Resource放到Session里
//func (controller *BaseController)SetSesstionRess(ress []models.Resource) {
//	var auths []string
//	for k, v := range ress {
//		beego.Info(k)
//		auths[k] = v.ActionPath
//	}
//	controller.SetSession(AUTHS, auths)
//}
//func (controller *BaseController) GetCurrentAdminId() int64 {
//	session := controller.GetSession(ADMIN_ID)
//	if adminId, ok := session.(int64); ok {
//		return adminId
//	}
//	return 0
//}
//
//func (controller *BaseController) GetCurrentAdminName() string {
//	session := controller.GetSession(ADMIN_NAME)
//	if adminName, ok := session.(string); ok {
//		return adminName
//	}
//	return ""
//}
const (
	OK_MSG = "操作成功"
	FAIL_MSG = "系统异常"
)
type JsonMsg struct {
	Ok   bool `json:"ok"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *BaseController) PrintOk() {
	this.printJson(true, OK_MSG, nil)
}

func (this *BaseController) PrintFail() {
	this.printJson(false, FAIL_MSG, nil)
}

func (this *BaseController) printJson(ok bool, msg string, data interface{}) {
	json := &JsonMsg{Ok: ok, Msg: msg, Data: data}
	this.Print(json)
}

func (this *BaseController) Print(data interface{}) {
	beego.Debug(fmt.Sprintf("%+v", data))
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *BaseController) PrintError(err error) {
	if err != nil {
		this.printJson(false, err.Error(), nil)
	} else {
		this.PrintOk()
	}
}

func (this *BaseController) PrintResult(data interface{}, err error) {
	if err != nil {
		this.printJson(false, err.Error(), nil)
	} else {
		if data != nil {
			this.Print(data)
		} else {
			this.PrintOk()
		}
	}
}

