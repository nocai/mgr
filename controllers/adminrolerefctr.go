package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"mgr/models/service/arrefser"
)

type AdminRoleRefController struct {
	BaseController
}

func (this *AdminRoleRefController) Post () {
	this.debugInput()

	adminId, _ := this.GetInt64("adminId", 0)
	roleIdsStr := this.GetStrings("roleIds[]",[]string{})

	beego.Info("adminId = ", adminId)
	beego.Info("roleIds = ", roleIdsStr)

	var roleIds []int64
	for _, roleId := range roleIdsStr {
		_roleId, _ := strconv.ParseInt(roleId,10,64)
		roleIds = append(roleIds,  _roleId)
	}
	_, err := arrefser.GrantRole(adminId, roleIds)
	this.PrintError(err)
}