package roleser

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/models"
	"mgr/models/service"
	"mgr/util/sqler"
)

func FindRoleByRoleDatagridKey(key *RoleDatagridKey) []models.Role {
	o := orm.NewOrm()

	sqler := sqler.New(key.Key)
	sqler.AppendSql(`select t.* from t_mgr_role t
		left join t_mgr_admin_role_ref tt on tt.role_id = t.id
		where tt.admin_id = ?`)
	sqler.SetAlias("t")
	var roles []models.Role

	affect, err := o.Raw(sqler.GetSql(), key.AdminId).QueryRows(&roles)
	if err != nil {
		panic(service.NewError(service.MsgQuery, err))
	}
	beego.Info("affect = ", affect)
	return roles
	//panic(service.NewError(service.MsgQuery, service.ErrArgument))
}
