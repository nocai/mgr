package role
import (
	"github.com/astaxie/beego/orm"
	"mgr/util"
	"github.com/astaxie/beego"
	"errors"
	"fmt"
	"time"
	"mgr/models"
	"mgr/models/service"
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
		return errors.New("删除失败")
	}

	beego.Debug("affected = %d", affected)
	return nil
}

// 取Role By Id
func GetRoleById(id int64) (*models.Role, error) {
	key := &models.RoleKey{Role:models.Role{Id:id}}
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

// 角色名存在
func isExistOfRole(role *models.Role) (bool, error) {
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
		return service.ErrInsert
	}
	beego.Debug("Id = %d", id)
	return nil
}

// 分页
func PageRole(key *models.RoleKey) (*util.Pager, error) {
	total, err := countRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, []models.Role{}), err
	}

	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, []models.Role{}), err
	}

	return util.NewPager(key.Key, total, roles), nil
}

func countRoleByKey(key *models.RoleKey) (int64, error) {
	o := orm.NewOrm()
	sqler := key.GetSqler()

	var total int64
	err := o.Raw(sqler.GetCountSql(), sqler.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return 0, service.ErrQuery
	}
	return total, nil
}

func FindRoleByKey(key *models.RoleKey) ([]models.Role, error) {
	o := orm.NewOrm()
	sqler := key.GetSqler()

	roles := []models.Role{}
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&roles)
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

// 取角色 By RoleName
func GetRoleByRoleName(roleName string) (*models.Role, error) {
	key := &models.RoleKey{Role:models.Role{RoleName:roleName}}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	if len(roles) == 0 {
		return nil, nil
	}
	if len(roles) > 1 {
		beego.Error(service.ErrDataDuplication, "roleName = ", roleName)
	}
	return &roles[0], nil
}
