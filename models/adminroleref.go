package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"mgr/util"
	"sync"
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

type AdminRoleRefKey struct {
	*util.Key
	AdminRoleRef
}

func (this *AdminRoleRefKey) getSelqer() *util.Sqler{
	sqler := &util.Sqler{Key:this.Key}
	sqler.AppendSql(`select * from t_mgr_admin_role_ref where 1 = 1`)

	if id := this.Id; id != 0 {
		sqler.AppendSql(" and id = ?")
		sqler.AppendArg(id)
	}
	if adminId := this.AdminId; adminId != 0 {
		sqler.AppendSql(" and admin_id = ?")
		sqler.AppendArg(adminId)
	}
	if roleId := this.RoleId; roleId != 0 {
		sqler.AppendSql(" and role_id = ?")
		sqler.AppendArg(roleId)
	}
	return sqler
}

func FindAdminRoleRefByKey(key *AdminRoleRefKey) ([]AdminRoleRef, error) {
	o := orm.NewOrm()
	sqler := key.getSelqer()

	var refs []AdminRoleRef
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&refs)
	if err != nil {
		beego.Error(err)
		return []AdminRoleRef{}, ErrQuery
	}
	beego.Debug("affected = ", affected)
	if affected == 0 {
		return []AdminRoleRef{}, nil
	}
	return refs, nil
}

func FindAdminByRoleId(roleId int64) ([]Admin, error) {
	key := &AdminRoleRefKey{AdminRoleRef:AdminRoleRef{RoleId:roleId}}
	refs, err := FindAdminRoleRefByKey(key)
	if err != nil {
		beego.Error(err)
		return []Admin{}, ErrQuery
	}
	if len(refs) == 0 {
		return []Admin{}, nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(refs))

	var result []Admin
	for _, ref := range refs {
		go func() {
			defer wg.Done()
			aKey := &AdminKey{Admin:Admin{Id:ref.AdminId}}
			admins, err := FindAdminByKey(aKey)
			if err != nil {
				beego.Error(err)
			}
			for _, admin := range admins {
				result = append(result, admin)
			}
		}()
	}
	return result, nil
}

func FindRoleByAdminId(adminId int64) ([]Role, error) {
	key := &AdminRoleRefKey{AdminRoleRef:AdminRoleRef{AdminId:adminId}}
	refs, err := FindAdminRoleRefByKey(key)
	if err != nil {
		beego.Error(err)
		return []Role{}, ErrQuery
	}
	if len(refs) == 0 {
		return []Role{}, nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(refs))

	var result []Role
	for _, ref := range refs {
		go func() {
			defer wg.Done()
			rKey := &RoleKey{Role:Role{Id:ref.RoleId}}
			roles, err := FindRoleByKey(rKey)
			if err != nil {
				beego.Error(err)
			}
			for _, role := range roles {
				result = append(result, role)
			}
		}()
	}
	return result, nil
}