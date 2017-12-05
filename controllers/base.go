package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"mgr/conf"
)

const (
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
panic(err)
		//this.PrintJson(err)

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

type ResultMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *BaseController) PrintJson(ret interface{}) {
	if ret == nil {
		this.Data["json"] = ResultMsg{Code: conf.CodeSuccess, Msg: conf.MsgSuccess}
	} else {
		switch ret.(type) {
		case error:
			this.Data["json"] = ResultMsg{Code: conf.CodeFail, Msg: ret.(error).Error()}
		case string:
			this.Data["json"] = ResultMsg{Code: conf.CodeSuccess, Msg: ret.(string)}
		default:
			this.Data["json"] = ResultMsg{Code: conf.CodeSuccess, Msg: conf.MsgSuccess, Data: ret}
		}
	}
	beego.Debug(fmt.Sprintf("%+v", this.Data["json"]))
	this.ServeJSON()
}


