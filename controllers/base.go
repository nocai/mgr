package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

const (
	OK  int = 200
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

func (this *BaseController) recoverPanic() {
	if err := recover(); err != nil {
		beego.Error("请求路径:", fmt.Sprintf("%#v", this.Ctx.Input.Params()))
		beego.Error("输入参数:", this.Input())
		beego.Error(err)
		switch err.(type) {
		case error:
			this.PrintError(err.(error))
		case string:
			this.PrintFail(err.(string))
		default:
			this.PrintData(err)
		}
	}
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
	OK_MSG   = "操作成功"
	FAIL_MSG = "系统异常"
)

type JsonMsg struct {
	Ok   bool        `json:"ok"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *BaseController) PrintOk(msg string) {
	if msg == "" {
		msg = OK_MSG
	}
	this.printJson(true, msg, nil)
}

func (this *BaseController) PrintFail(msg string) {
	if msg == "" {
		msg = FAIL_MSG
	}
	this.printJson(false, msg, nil)
}

func (this *BaseController) printJson(ok bool, msg string, data interface{}) {
	json := &JsonMsg{Ok: ok, Msg: msg, Data: data}
	this.doPrint(json)
}

func (this *BaseController) doPrint(data interface{}) {
	beego.Debug(fmt.Sprintf("%+v", data))
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *BaseController) PrintData(data interface{}) {
	this.doPrint(data)
}

func (this *BaseController) PrintError(err error) {
	if err != nil {
		this.PrintFail(err.Error())
	} else {
		this.PrintOk("")
	}
}

func (this *BaseController) Print(data interface{}, err error) {
	if err != nil {
		this.PrintFail(err.Error())
		return
	}
	if data != nil {
		this.doPrint(data)
		return
	}
	panic("data and err are all nil")
}
