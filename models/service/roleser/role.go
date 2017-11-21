package roleser

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/models"
	"mgr/models/service"
	"mgr/models/service/arrefser"
	"mgr/util/key"
	"mgr/util/pager"
	"sync"
	"time"
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
	affected, err := o.Delete(&models.Role{Id: id})
	if err != nil {
		beego.Error(err)
		return service.ErrDelete
	}

	beego.Debug("DeleteRoleById: affected = ", affected)
	return nil
}

// 取Role By Id
func GetRoleById(id int64) (*models.Role, error) {
	key := &models.RoleKey{Role: &models.Role{Id: id}}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if len(roles) == 0 {
		return nil, orm.ErrNoRows
	}
	if len(roles) > 1 {
		beego.Error(service.ErrDataDuplication)
	}
	return &roles[0], nil
}

// isExistOfRole.
// When the error != nil, the bool is invalid.
func isExistOfRole(role *models.Role) (bool, error) {
	roleId := role.Id
	// 设置Id = 0，方便查询
	role.Id = 0
	key := &models.RoleKey{Role: role}
	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return false, err
	}

	for _, _role := range roles {
		beego.Debug("the exist role = ", _role)
		if _role.Id != roleId {
			beego.Debug(fmt.Sprintf("the role is exist: role = %v", role))
			beego.Debug(fmt.Sprintf("the role in db: role = %v", _role))
			return true, nil
		}
	}
	// 将Id设置回来，不然role的数据不对
	role.Id = roleId
	beego.Debug(fmt.Sprintf("the role not exist: role = %v", role))
	return false, nil
}

// 更新
func UpdateRole(role *models.Role) error {
	ormer := orm.NewOrm()
	exist, err := isExistOfRole(&models.Role{
		Id:       role.Id,
		RoleName: role.RoleName,
	})
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
	beego.Debug("InsertRole: Id = ", id)
	role.Id = id
	return nil
}

// 分页
func PageRole(key *models.RoleKey) (*pager.Pager, error) {
	total, err := countRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return pager.New(key.Key, 0, []models.Role{}), err
	}

	roles, err := FindRoleByKey(key)
	if err != nil {
		beego.Error(err)
		return pager.New(key.Key, 0, []models.Role{}), err
	}

	return pager.New(key.Key, total, roles), nil
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

	key := &models.RoleKey{Role: &models.Role{RoleName: roleName}}
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

func FindRoleByAdminId(adminId int64) ([]models.Role, error) {
	arrefKey := &models.AdminRoleRefKey{
		AdminRoleRef: &models.AdminRoleRef{AdminId: adminId},
		Key:          &key.Key{Sort: []string{"id"}, Order: []string{"desc"}},
	}

	refs, err := arrefser.FindAdminRoleRefByKey(arrefKey)
	beego.Error(refs)
	if err != nil {
		beego.Error(err)
		return []models.Role{}, service.NewError(service.MsgQuery, err)
	}
	if len(refs) == 0 {
		return []models.Role{}, nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(refs))

	var result []models.Role
	for i := range refs {
		go func(index int) {
			defer wg.Done()
			rKey := &models.RoleKey{Role: &models.Role{Id: refs[index].RoleId}, Key: &key.Key{Sort: []string{"id"}, Order: []string{"desc"}}}
			roles, err := FindRoleByKey(rKey)
			if err != nil {
				beego.Error(err)
			}
			for _, role := range roles {
				result = append(result, role)
			}
		}(i)
	}
	wg.Wait()
	beego.Error(result)
	return result, nil
}
