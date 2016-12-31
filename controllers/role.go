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

func (ctr *RoleController) Post() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	roleName := ctr.GetString("role_name")

	now := time.Now()
	role := &models.Role{Id:id, RoleName:roleName, CreateTime:now, UpdateTime:now}
	if id == 0 {
		err := models.InsertRole(role)
		if err != nil {
			beego.Error(err)
			ctr.PrintErrorMsg(err.Error())
			return
		}
		ctr.PrintOk()
	} else {
		r, err := models.GetRoleById(role.Id)
		if err != nil {
			beego.Error(err)
			ctr.PrintErrorMsg(err.Error())
			return
		}

		r.RoleName = role.RoleName
		r.UpdateTime = time.Now()
		err = models.UpdateRole(r)
		if err != nil {
			beego.Error(err)
			ctr.PrintErrorMsg(err.Error())
			return
		}
		ctr.PrintOk()
	}
}

func (ctr *RoleController) Get() {
	beego.Debug(ctr.Input())

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")

	roleName := ctr.GetString("role_name")
	data := make(map[string]interface{}, 1)
	data["roleName"] = roleName

	pagerKey := util.NewPagerKey(page, rows, data, sort, order)
	pager, err := models.PageRole(pagerKey)
	if err != nil {
		beego.Error(err)
		ctr.Print(util.NewPager(pagerKey, 0, make([]interface{}, 0)))
		return
	}

	ctr.Print(pager.Pagination)
}