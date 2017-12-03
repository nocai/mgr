package adminser

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/models"
	"mgr/models/service"
	"mgr/models/service/arrefser"
	"mgr/util/pager"
	"sync"
	"time"
	"mgr/conf"
)

var (
	ErrUsernameNotExist   = errors.New("用户名不存在")
	ErrUsernameExist      = errors.New("用户名存在")
	ErrPasswordNotMatched = errors.New("密码错误")
	ErrNotSysAdmin        = errors.New("对不起,您还不是系统管理员")
)

func CountAdminByKey(key *models.AdminKey) (int64, error) {
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

func FindAdminByKey(key *models.AdminKey) ([]models.Admin, error) {
	o := orm.NewOrm()
	sqler := key.NewSqler()

	var admins []models.Admin
	affected, err := o.Raw(sqler.GetSqlAndArgs()).QueryRows(&admins)
	if err != nil {
		beego.Error(err)
		return admins, err
	}
	beego.Debug(fmt.Sprintf("affected = %d", affected))
	if affected == 0 {
		return []models.Admin{}, nil
	}
	return admins, nil
}

func PageAdmin(key *models.AdminKey) (*pager.Pager, error) {
	total, err := CountAdminByKey(key)
	if err != nil {
		beego.Error(err)
		return pager.New(key.Key, 0, []models.Admin{}), service.ErrQuery
	}

	admins, err := FindAdminByKey(key)
	if err != nil {
		beego.Error(err)
		return pager.New(key.Key, 0, []models.Admin{}), service.ErrQuery
	}
	return pager.New(key.Key, total, admins), nil
}

func InsertAdmin(admin *models.Admin) error {
	if admin.UserId == 0 {
		beego.Error("管理员的UserId必须填")
		return service.ErrArgument
	}

	now := time.Now()
	if admin.CreateTime.IsZero() {
		admin.CreateTime = now
	}
	if admin.UpdateTime.IsZero() {
		admin.UpdateTime = now
	}

	o := orm.NewOrm()
	id, err := o.Insert(admin)
	if err != nil {
		beego.Error(err)
		return service.ErrInsert
	}
	beego.Debug(fmt.Sprintf("Add Admin sucess: id = %v", id))
	return nil
}

func DeleteAdminById(id int64) {
	if id == 0 {
		panic(service.ErrArgument)
	}

	key := &models.AdminKey{Admin: &models.Admin{Id: id}}
	admins, err := FindAdminByKey(key)
	if err != nil {
		panic(service.NewError(conf.MsgQuery, err))
	}
	if len(admins) == 0 {
		return
	}
	if len(admins) > 1 {
		panic(service.ErrDataDuplication)
	}

	o := orm.NewOrm()
	o.Begin()

	affected, err := o.Delete(&models.User{Id: admins[0].UserId})
	if err != nil {
		o.Rollback()
		panic(service.NewError(conf.MsgQuery, err))
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	affected, err = o.Delete(&models.Admin{Id: id})
	if err != nil {
		o.Rollback()
		panic(service.NewError(conf.MsgDelete, err))
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	affected, err = o.Delete(&models.AdminRoleRef{AdminId: id}, "admin_id")
	if err != nil {
		o.Rollback()
		panic(service.NewError(conf.MsgDelete, err))
	}
	beego.Debug("affected = ", affected)

	o.Commit()
}

func GetAdminById(id int64) (*models.Admin, error) {
	adminKey := &models.AdminKey{
		Admin: &models.Admin{
			Id: id,
		},
	}
	adminSlice, err := FindAdminByKey(adminKey)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if len(adminSlice) == 0 {
		beego.Error(orm.ErrNoRows)
		return nil, orm.ErrNoRows
	} else if len(adminSlice) == 1 {
		return &adminSlice[0], nil
	} else {
		beego.Error(fmt.Sprintf("data duplication: id = %d", id))
		return nil, service.ErrDataDuplication
	}
}

// 取所有的Admin valid
// key:models.ValidEnum value:string
func FindAdminValids() []map[string]interface{} {
	allMap := make(map[string]interface{})
	allMap["value"] = models.ValidAll
	allMap["text"] = "全部"

	invalidMap := make(map[string]interface{})
	invalidMap["value"] = models.Invalid
	invalidMap["text"] = "无效"

	validMap := make(map[string]interface{})
	validMap["value"] = models.Valid
	validMap["text"] = "有效"

	return []map[string]interface{}{
		allMap, invalidMap, validMap,
	}
}

func UpdateAdmin(admin *models.Admin) error {
	if admin == nil {
		return service.ErrArgument
	}
	o := orm.NewOrm()
	num, err := o.Update(admin)
	if err != nil {
		beego.Error(err)
		return err
	}
	beego.Info("<UpdateAdmin>: num = ", num)
	return nil
}

func FindAdminByRoleId(roleId int64) ([]models.Admin, error) {
	key := &models.AdminRoleRefKey{AdminRoleRef: &models.AdminRoleRef{RoleId: roleId}}
	refs, err := arrefser.FindAdminRoleRefByKey(key)
	if err != nil {
		beego.Error(err)
		return []models.Admin{}, service.NewError(conf.MsgQuery, err)
	}
	if len(refs) == 0 {
		return []models.Admin{}, nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(refs))

	var result []models.Admin
	for _, ref := range refs {
		go func() {
			defer wg.Done()
			aKey := &models.AdminKey{Admin: &models.Admin{Id: ref.AdminId}}
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
