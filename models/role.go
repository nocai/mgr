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
	ModelBase

	Id         int64  `json:"id"`
	RoleName   string `orm:"unique" json:"role_name"`
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
func PageRole(key *RoleKey) (*util.Pager, error) {
	total, err := countRoleByKey(key)
	if err != nil {
		return util.NewPager(key.Key, 0, []Role{}), err
	}

	roles, err := FindRoleByKey(key)
	if err != nil {
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

	var roles []Role
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&roles)
	if err != nil {
		beego.Error(err)
		return roles, ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return roles, nil
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