package controllers

import (
	"mgr/models"
	"github.com/astaxie/beego"
	"time"
	"mgr/conf"
	"fmt"
	"mgr/models/service/roleser"
	"mgr/util/key"
)

type RoleController struct {
	BaseController
}

// 删除
func (ctr *RoleController) Delete() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	beego.Error(id)

	err := roleser.DeleteRoleById(id)
	ctr.PrintError(err)
}

// 添加 修改
func (ctr *RoleController) Post() {
	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug(fmt.Sprintf("id = %v", id))
	beego.Debug(ctr.Input())

	roleName := ctr.GetString("role_name")
	if id == 0 {
		beego.Error("add")
		ctr.PrintError(addRole(roleName))
	} else {
		beego.Error("update")
		ctr.PrintError(updateRole(id, roleName))
	}
}

func updateRole(id int64, roleName string) error {
	r, err := roleser.GetRoleById(id)
	if err != nil {
		return err
	}
	r.RoleName = roleName
	r.UpdateTime = time.Now()
	err = roleser.UpdateRole(r)
	return err
}

func addRole(roleName string) error {
	r := &models.Role{RoleName:roleName}
	return roleser.InsertRole(r);
}

func (ctr *RoleController) Get() {
	ctr.debugInput()

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	roleName := ctr.GetString("role_name")

	key := key.New(page, rows, []string{sort}, []string{order}, true)
	r := &models.Role{RoleName:roleName}
	roleKey := &models.RoleKey{Key:key, Role:r}
	pager, err := roleser.PageRole(roleKey)
	if err != nil {
		beego.Error(err)
	}
	ctr.Print(pager.Pagination)
}