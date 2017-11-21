package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"mgr/models/service/arrefser"
	"strconv"
)

type AdminRoleRefController struct {
	BaseController
}

func (this *AdminRoleRefController) Post() {
	adminId, _ := this.GetInt64("adminId", 0)
	roleIdsStr := this.GetStrings("roleIds[]", []string{})

	beego.Info("adminId = ", adminId)
	beego.Info("roleIds = ", roleIdsStr)

	var roleIds []int64
	for _, roleId := range roleIdsStr {
		_roleId, _ := strconv.ParseInt(roleId, 10, 64)
		roleIds = append(roleIds, _roleId)
	}
	arIds := arrefser.GrantRole(adminId, roleIds)
	beego.Debug(fmt.Sprintf("%#v", arIds))
	this.PrintOk("授权成功")
}
