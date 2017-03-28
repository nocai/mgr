package models

import (
	"github.com/astaxie/beego/orm"
	"mgr/util"
	"github.com/astaxie/beego"
	"errors"
	"fmt"
	"time"
)

var (
	ErrRoleNameExist = errors.New("角色名称存在")
)

type Role struct {
	ModelBase

	Id       int64  `json:"id"`
	RoleName string `orm:"unique" json:"role_name"`
}

// 多字段索引
func (role *Role) TableIndex() [][]string {
	return [][]string{
		[]string{"RoleName"},
	}
}

type RoleKey struct {
	*util.Key

	Role
	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *RoleKey) getSqler() *util.Sqler {
	sqler := &util.Sqler{Key:this.Key}

	sqler.AppendSql("select * from t_mgr_role as tmr where 1 = 1")
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tmr.id = ?")
		sqler.AppendArg(id)
	}
	if roleName := this.RoleName; roleName != "" {
		sqler.AppendSql(" and tmr.role_name like ?")
		sqler.AppendArg("%" + roleName + "%")
	}
	return sqler
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
	key := &RoleKey{Role:Role{Id:id}}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	if len(roles) > 1 {
		beego.Error(ErrDataDuplication)
	}
	return &roles[0], nil
}

// 角色名存在
func isExistOfRole(role *Role) (bool, error) {
	r, err := GetRoleByRoleName(role.RoleName)
	if err != nil {
		beego.Error(err)
		return false, err
	}

	if r == nil {
		return false, nil
	}
	if r.Id != role.Id {
		beego.Debug(fmt.Sprintf("角色名存在。id = %v, roleName = %v", r.Id, r.RoleName))
		return true, nil
	}
	return false, nil
}



// 更新
func UpdateRole(role *Role) error {
	ormer := orm.NewOrm()
	exist, err := isExistOfRole(role);
	if err != nil {
		beego.Error(err)
		return err
	}
	if exist {
		return ErrRoleNameExist
	}

	affected, err := ormer.Update(role)
	if err != nil {
		beego.Error(err)
		return ErrUpdate
	}
	beego.Debug("<UpdateRole> affected = %v", affected)
	return nil
}

// 添加
func InsertRole(role *Role) error {
	ormer := orm.NewOrm()

	exist, err := isExistOfRole(role)
	if err != nil {
		beego.Error(err)
		return err
	}
	if exist {
		return ErrRoleNameExist
	}

	if role.CreateTime.IsZero() {
		role.CreateTime = time.Now()
	}
	if role.UpdateTime.IsZero() {
		role.UpdateTime = time.Now()
	}

	id, err := ormer.Insert(role)
	if err != nil {
		beego.Error(err)
		return ErrInsert
	}
	beego.Debug("Id = %d", id)
	return nil
}

// 分页
func PageRole(key *RoleKey) (*util.Pager, error) {
	total, err := countRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, []Role{}), err
	}

	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, []Role{}), err
	}

	return util.NewPager(key.Key, total, roles), nil
}

func countRoleByKey(key *RoleKey) (int64, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var total int64
	err := o.Raw(sqler.GetCountSql(), sqler.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return 0, ErrQuery
	}
	return total, nil
}

func FindRoleByKey(key *RoleKey) ([]Role, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	roles := []Role{}
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&roles)
	if err != nil {
		beego.Error(err)
		return roles, ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	if affected == 0 {
		beego.Debug(orm.ErrNoRows)
		return []Role{}, nil
	}
	return roles, nil
}

// 取角色 By RoleName
func GetRoleByRoleName(roleName string) (*Role, error) {
	key := &RoleKey{Role:Role{RoleName:roleName}}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	if len(roles) == 0 {
		return nil, nil
	}
	if len(roles) > 1 {
		beego.Error(ErrDataDuplication, "roleName = ", roleName)
	}
	return &roles[0], nil
}
