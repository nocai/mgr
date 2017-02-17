package models

import (
	"github.com/astaxie/beego/orm"
	"mgr/util"
	"github.com/astaxie/beego"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Role struct {
	Id         int64  `json:"id"`
	RoleName   string `orm:"unique" json:"role_name"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}



// 多字段索引
func (role *Role) TableIndex() [][]string {
	return [][]string{
		[]string{"RoleName"},
	}
}
// 删除
func DeleteRoleById(id int64) error {
	if id <= 0 {
		beego.Error("Id = %d", id)
		return nil
	}

	o := orm.NewOrm()
	affected, err := o.Delete(&Role{Id:id})
	if err != nil {
		beego.Error(err)
		return errors.New("删除失败")
	}

	beego.Debug("affected = %d", affected)
	return nil
}

// 取Role By Id
func GetRoleById(id int64) (*Role, error) {
	ormer := orm.NewOrm()

	role := &Role{Id:id}
	err := ormer.Read(role)
	if err != nil {
		beego.Error(fmt.Sprintf("查询失败.Id = %v", id), err)
		return nil, errors.New("查询失败")
	}
	return role, nil
}

// 角色名存在
func IsExist(role *Role) bool {
	if roleName := role.RoleName; roleName != "" {
		r, err := GetRoleByRoleName(roleName)
		if err == nil {
			if r.Id != role.Id {
				beego.Info("角色名存在:" + roleName)
				return true
			}
		}
	}
	return false
}

// 取角色 By RoleName
func GetRoleByRoleName(roleName string) (*Role, error) {
	o := orm.NewOrm()

	role := &Role{RoleName:roleName}
	err := o.Read(role, "RoleName")

	if err != nil {
		beego.Error(err)
		return nil, errors.New("查询失败")
	}
	return role, nil
}

// 更新
func UpdateRole(role *Role) error {
	ormer := orm.NewOrm()

	if role.Id == 0 {
		panic(fmt.Errorf("role.Id = %d", role.Id))
	}
	if IsExist(role) {
		return errors.New("角色名称存在")
	}

	affected, err := ormer.Update(role)
	if err != nil {
		beego.Error(err)
		return errors.New("更新失败")
	}
	beego.Debug("<UpdateRole> affected = %v", affected)
	return nil
}

// 添加
func InsertRole(role *Role) error {
	ormer := orm.NewOrm()

	if IsExist(role) {
		return errors.New("角色名称存在")
	}

	id, err := ormer.Insert(role)
	if err != nil {
		beego.Error(err)
		return errors.New("添加失败")
	}
	beego.Debug("Id = %d", id)
	return nil
}

// 分页
func PageRole(key *util.PagerKey) (*util.Pager, error) {
	key.AppendSql("select * from t_mgr_role as tmr where 1 = 1")

	if roleName, ok := key.Data["roleName"].(string); ok && roleName != "" {
		key.AppendSql(" and tmr.role_name like ?")
		key.AppendArg("%" + roleName + "%")
	}

	ormer := orm.NewOrm()

	var total int64
	err := ormer.Raw(key.GetCountSql(), key.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return nil, errors.New("分页失败")
	}

	var roles []Role
	affected, err := ormer.Raw(key.GetSql(), key.GetArgs()).QueryRows(&roles)
	if err != nil {
		beego.Error(err)
		return nil, errors.New("分页失败")
	}

	beego.Debug("affected = %d", affected)
	return util.NewPager(key, total, roles), nil
}

func FindRolesByUserId(userId int64) (*[]Role, error) {
	refs, err := FindAdminRoleRefByAdminId(userId)
	if err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}

	wg := sync.WaitGroup{}
	wg.Add(len(*refs))

	var roles []Role
	for _, ref := range *refs {
		ref := ref
		go func() {
			defer wg.Done()

			role, err := GetRoleById(ref.RoleId)
			if err != nil {
				beego.Error(err)
				return
			}
			roles = append(roles, *role)
		}()
	}

	wg.Wait()

	return &roles, nil
}