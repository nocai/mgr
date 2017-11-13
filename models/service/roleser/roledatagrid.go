package roleser

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/models"
	"mgr/models/service"
)

func FindRoleByRoleDatagridKey(key *RoleDatagridKey) []models.Role {
	o := orm.NewOrm()

	var roles []models.Role
	var query = `select t.* from t_mgr_role t
		left join t_mgr_admin_role_ref tt on tt.role_id = t.id
		where tt.admin_id = ?`
	affect, err := o.Raw(query+key.GetOrderBySql("t")+key.GetLimitSql(), key.AdminId).QueryRows(&roles)
	if err != nil {
		panic(service.NewError(service.MsgQuery, err))
	}
	beego.Info("affect = ", affect)
	return roles
	//panic(service.NewError(service.MsgQuery, service.ErrArgument))
}
