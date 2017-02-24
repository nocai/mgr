package controllers

import (
	"mgr/util"
	"mgr/models"
	"github.com/astaxie/beego"
	"time"
	"mgr/conf"
)

type RoleController struct {
	BaseController
}

// 删除
func (ctr *RoleController) Delete() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	beego.Error(id)

	err := models.DeleteRoleById(id)
	if err != nil {
		beego.Error(err)
		ctr.PrintErrorMsg(err.Error())
		return
	}
	ctr.PrintOk()
}

// 添加 修改
func (ctr *RoleController) Post() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	roleName := ctr.GetString("role_name")

	var err error
	if id == 0 {
		err = addRole(roleName)
	} else {
		err = updateRole(id, roleName)
	}

	if err != nil {
		beego.Error(err)
		ctr.PrintErrorMsg(err.Error())
		return
	}
	ctr.PrintOk()
}

func updateRole(id int64, roleName string) error {
	role, err := models.GetRoleById(id)
	if err != nil {
		return err
	}

	role.RoleName = roleName
	role.UpdateTime = time.Now()
	err = models.UpdateRole(role)
	return err
}

func addRole(roleName string) error {
	role := &models.Role{RoleName:roleName}
	return models.InsertRole(role);
}

func (ctr *RoleController) Get() {
	ctr.debugInput()

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	roleName := ctr.GetString("role_name")
	data := make(map[string]interface{}, 1)
	data["roleName"] = roleName

	key := util.NewKey(page, rows, []string{sort}, []string{order}, true)
	pager, err := models.PageRole(&models.RoleKey{Key:key})
	if err != nil {
		beego.Error(err)
	}
	ctr.Print(pager.Pagination)
}