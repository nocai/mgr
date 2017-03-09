package models

import (
	"mgr/util"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"strconv"
)

const (
	// 开启
	PrivilegeOperation_Enable = iota
	// 禁用
	PrivilegeOperation_Disabled
)

const (
	// 角色
	PrivilegeMaster_Role = iota
	// 用户
	PrivilegeMaster_Admin
)

const (
	// 菜单
	PrivilegeAccess_Menu = ResType_Menu
	// 按钮
	PrivilegeAccess_Button = ResType_Button
)

// 某某主体 在 某某领域 有 某某权限
type Privilege struct {
	ModelBase
	Id                   int64

	PrivilegeMaster      string
	PrivilegeMasterValue int64

	PrivilegeAccess      string
	PrivilegeAccessValue int64

	PrivilegeOperation   string
}

type PrivilegeKey struct {
	*util.Key

	Privilege
}

func (this *PrivilegeKey) getSqler() *util.Sqler {
	sqler := &util.Sqler{Key:this.Key}

	sqler.AppendSql(`select * from t_mgr_privilege as tmp where 1 = 1`)
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tmp.id = ?")
		sqler.AppendArg(id)
	}
	if privilegeMaster := this.PrivilegeMaster; privilegeMaster != "" {
		sqler.AppendSql(" and tmp.privilege_master = ?")
		sqler.AppendArg(privilegeMaster)
	}
	if privilegeMasterValue := this.PrivilegeMasterValue; privilegeMasterValue != 0 {
		sqler.AppendSql(" and tmp.privilege_master_value = ?")
		sqler.AppendArg(privilegeMasterValue)
	}
	if privilegeAccess := this.PrivilegeAccess; privilegeAccess != "" {
		sqler.AppendSql(" and tmp.privilege_access = ?")
		sqler.AppendArg(privilegeAccess)
	}
	if privilegeAccessValue := this.PrivilegeAccessValue; privilegeAccessValue != 0 {
		sqler.AppendSql(" and tmp.privilege_access_value = ?")
		sqler.AppendArg(privilegeAccessValue)
	}
	if privilegeOperation := this.PrivilegeOperation; privilegeOperation != "" {
		sqler.AppendSql(" and tmp.privilege_operation = ?")
		sqler.AppendArg(privilegeOperation)
	}
	return sqler
}

func FindPrivilegeByKey(key *PrivilegeKey) ([]Privilege, error) {
	sqler := key.getSqler()

	o := orm.NewOrm()
	var privileges []Privilege
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&privileges)
	if err != nil {
		beego.Error(err)
		return []Privilege{}, ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	if affected == 0 {
		return []Privilege{}, nil
	}
	return privileges, nil
}

func FindResByAdminId(adminId int64) ([]Res, error) {
	ress, err := findResBelongRole(adminId)
	if err != nil {
		beego.Error(err)
		return []Res{}, ErrQuery
	}

	ress2, err := findResBelongAdmin(adminId)
	if err != nil {
		beego.Error(err)
		return []Res{}, ErrQuery
	}
	ress = append(ress, ress2...)

	// 去重复数据
	var result []Res
	for _, res := range ress {
		if !func() bool {
			for _, v := range result {
				if v.Id == res.Id {
					return true
				}
			}
			return false
		}() {
			result = append(result, res)
		}
	}
	return result, nil
}

// 取adminId所有权限
func findResBelongAdmin(adminId int64) ([]Res, error) {

	key := &PrivilegeKey{Privilege:Privilege{PrivilegeMaster:strconv.Itoa(PrivilegeMaster_Admin), PrivilegeMasterValue:adminId}}
	privileges, err := FindPrivilegeByKey(key)
	if err != nil {
		beego.Error(err)
		return []Res{}, ErrQuery
	}

	result := []Res{}
	for _, privilege := range privileges {
		key2 := &ResKey{Res:Res{Id:privilege.PrivilegeMasterValue}}
		ress, err := FindResByKey(key2)
		if err != nil {
			beego.Error(err)
			return result, ErrQuery
		}
		result = append(result, ress...)
	}
	return result, nil
}

// 取adminId所属角色的权限
func findResBelongRole(adminId int64) ([]Res, error) {
	roles, err := FindRoleByAdminId(adminId)
	if err != nil {
		beego.Error(err)
		return []Res{}, ErrQuery
	}

	result := []Res{}
	for _, role := range roles {
		key := &PrivilegeKey{Privilege:Privilege{PrivilegeMaster:strconv.Itoa(PrivilegeMaster_Role), PrivilegeAccessValue:role.Id}}
		privileges, err := FindPrivilegeByKey(key)
		if err != nil {
			beego.Error(err)
			return []Res{}, ErrQuery
		}
		for _, privilege := range privileges {
			key2 := &ResKey{Res:Res{Id:privilege.PrivilegeMasterValue}}
			ress, err := FindResByKey(key2)
			if err != nil {
				beego.Error(err)
				return []Res{}, ErrQuery
			}
			result = append(result, ress...)
		}
	}
	return result, nil
}

