package controllers

import (
	"github.com/astaxie/beego"
	"mgr/models"
	"mgr/models/service/roleser"
	"mgr/util/key"
)

type RoleDatagridController struct {
	BaseController
}

type RoleDatagrid struct {
	*models.Role
	Checked bool `json:"checked"`
}

func (this *RoleDatagridController) Get() {
	defer this.recoverPanic()
	sort := this.GetString("sort")
	order := this.GetString("order")
	adminId, _ := this.GetInt64(":adminId")

	key := key.New(0, 0, []string{sort}, []string{order})
	roleKey := &models.RoleKey{Key: key}
	pager, err := roleser.PageRole(roleKey)
	if err != nil {
		beego.Error(err)
	}

	rdKey := &roleser.RoleDatagridKey{Key: key, AdminId: adminId}
	roles := roleser.FindRoleByRoleDatagridKey(rdKey)

	var rds []RoleDatagrid
	if allRoles, ok := pager.Pagination.PageList.([]models.Role); ok {
		for i := 0; i < len(allRoles); i++ {
			rds = append(rds, RoleDatagrid{
				Role:    &allRoles[i],
				Checked: contains(roles, &allRoles[i]),
			})
		}

	}
	pager.Pagination.PageList = rds
	this.PrintData(pager)
}

func contains(roles []models.Role, role *models.Role) bool {
	for _, r := range roles {
		if r.Id == role.Id {
			return true
		}
	}
	return false
}
