package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
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

type JsonMsg struct {
	Ok   bool `json:"ok"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *BaseController) PrintOk() {
	this.PrintOkMsg("操作成功")
}

func (this *BaseController) PrintOkMsg(msg string) {
	this.PrintOkMsgData(msg, nil)
}

func (this *BaseController) PrintOkMsgData(msg string, data interface{}) {
	this.printJson(true, msg, data)
}

func (this *BaseController) PrintError() {
	this.PrintErrorMsg("系统异常")
}

func (this *BaseController) PrintErrorMsg(msg string) {
	this.PrintErrorMsgData(msg, nil)
}

func (this *BaseController) PrintErrorMsgData(msg string, data interface{}) {
	this.printJson(false, msg, data)
}

func (this *BaseController) printJson(ok bool, msg string, data interface{}) {
	json := &JsonMsg{Ok: ok, Msg: msg, Data: data}
	this.Print(json)
}

func (this *BaseController) Print(data interface{}) {
	beego.Debug(data)
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *BaseController) PrintErrorMsgValid(r *validation.Result) {
	this.PrintErrorMsg(r.Error.Key + ", " + r.Error.Message)
}
