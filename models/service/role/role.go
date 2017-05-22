package role

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"errors"
	"fmt"
	"time"
	"mgr/models"
	"mgr/models/service"
	"mgr/util/pager"
)

var (
	ErrRoleNameExist = errors.New("角色名称存在")
)

// 删除
func DeleteRoleById(id int64) error {
	if id <= 0 {
		beego.Error("Id = %d", id)
		return nil
	}

	o := orm.NewOrm()
	affected, err := o.Delete(&models.Role{Id:id})
	if err != nil {
		beego.Error(err)
		return service.ErrDelete
	}

	beego.Debug("affected = %d", affected)
	return nil
}

// 取Role By Id
func GetRoleById(id int64) (*models.Role, error) {
	key := &models.RoleKey{Role:&models.Role{Id:id}}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	if len(roles) > 1 {
		beego.Error(service.ErrDataDuplication)
	}
	return &roles[0], nil
}

// isExistOfRole.
// When the error != nil, the bool is invalid.
func isExistOfRole(role *models.Role) (bool, error) {
	beego.Info(fmt.Sprintf("检查角色名是否存在:role = %v", role))
	if role.RoleName == "" {
		beego.Debug(service.ErrArgument, fmt.Sprintf("role.RoleName = %s", role.RoleName))
		return false, service.ErrArgument
	}

	r, err := GetRoleByRoleName(role.RoleName)
	if err != nil {
		if err == orm.ErrNoRows {
			return false, nil
		} else {
			beego.Error(err)
			return false, err
		}
	}

	if r.Id != role.Id {
		beego.Debug(fmt.Sprintf("the role name is exist. id = %v, roleName = %v", r.Id, r.RoleName))
		return true, nil
	}
	return false, nil
}

// 更新
func UpdateRole(role *models.Role) error {
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
		return service.ErrUpdate
	}
	beego.Debug("<UpdateRole> affected = %v", affected)
	return nil
}

// 添加
func InsertRole(role *models.Role) error {
	// check the role name
	exist, err := isExistOfRole(role)
	if err != nil {
		beego.Error(err)
		return err
	}

	if exist {
		return ErrRoleNameExist
	}

	// init the create time and last update time
	now := time.Now()
	if role.CreateTime.IsZero() {
		role.CreateTime = now
	}
	if role.UpdateTime.IsZero() {
		role.UpdateTime = now
	}

	// insert into the database
	ormer := orm.NewOrm()
	id, err := ormer.Insert(role)
	if err != nil {
		beego.Error(err)
		return service.ErrInsert
	}
	beego.Debug("Id = %d", id)
	return nil
}

// 分页
func PageRole(key *models.RoleKey) (*pager.Pager, error) {
	total, err := countRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return pager.New(key, 0, []models.Role{}), err
	}

	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return pager.New(key, 0, []models.Role{}), err
	}

	return pager.New(key, total, roles), nil
}

func countRoleByKey(key *models.RoleKey) (int64, error) {
	o := orm.NewOrm()
	sqler := key.NewSqler()

	var total int64
	err := o.Raw(sqler.GetCountSqlAndArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return 0, service.ErrQuery
	}
	return total, nil
}

// Query roles by key from the database.
// The key:see models.RoleKey
// If no row in he database, the method will return empty slice and nil error
func FindRoleByKey(key *models.RoleKey) ([]models.Role, error) {
	o := orm.NewOrm()
	sqler := key.NewSqler()

	roles := []models.Role{}
	affected, err := o.Raw(sqler.GetSqlAndArgs()).QueryRows(&roles)
	if err != nil {
		beego.Error(err)
		return roles, service.ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	if affected == 0 {
		beego.Debug(orm.ErrNoRows)
		return []models.Role{}, nil
	}
	return roles, nil
}

// GetRoleByRoleName.
// If has no rows, will return nil and error:orm.ErrNoRows
func GetRoleByRoleName(roleName string) (*models.Role, error) {
	beego.Info(fmt.Sprintf("GetRoleByRoleName.roleName = %s", roleName))

	key := &models.RoleKey{Role:&models.Role{RoleName:roleName}}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	if len(roles) == 0 {
		return nil, orm.ErrNoRows
	} else if len(roles) > 1 {
		beego.Error(service.ErrDataDuplication, "roleName = ", roleName)
	}
	return &roles[0], nil
}
