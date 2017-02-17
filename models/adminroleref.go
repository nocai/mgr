package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
)

type AdminRoleRef struct {
	Id      int64
	AdminId int64
	RoleId  int64
}

func (ref *AdminRoleRef) TableIndex() [][]string {
	return [][] string{
		[] string{"AdminId"},
		[] string{"RoleId"},
	}
}
func FindAdminRoleRefByAdminId(adminId int64) (*[]AdminRoleRef, error) {
	if adminId == 0 {
		return nil, ErrArgument
	}

	o := orm.NewOrm()

	var refs []AdminRoleRef
	affected, err := o.Raw("select * from t_mgr_admin_role_ref where admin_id = ?", adminId).QueryRows(&refs)
	if err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}

	beego.Debug(fmt.Sprintf("affected = %d", affected))
	return &refs, nil
}

func FindAdminRoleRefByRoleId(roleId int64) (*[]AdminRoleRef, error) {
	if roleId == 0 {
		return nil, ErrArgument

	}
	o := orm.NewOrm()

	var refs *[]AdminRoleRef
	affected, err := o.Raw("select * from t_mgr_admin_role_ref where role_id = ?", roleId).QueryRows(refs)
	if err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}

	beego.Debug(fmt.Sprintf("affected = %d", affected))
	return refs, nil
}