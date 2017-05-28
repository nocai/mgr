package adminser

import (
	"mgr/models/service"
	"fmt"
	"github.com/astaxie/beego"
	"mgr/models"
	"github.com/astaxie/beego/orm"
	"time"
	"mgr/models/service/userser"
)
//
//import (
//	"github.com/astaxie/beego"
//	"github.com/astaxie/beego/orm"
//	"fmt"
//	"mgr/models"
//	"mgr/util/pager"
//	"mgr/models/service"
//	"time"
//)
//
//func Login(username, password string) (*models.AdminVo, error) {
//	user, err := GetUserByUsername(username)
//	if err != nil {
//		beego.Debug(err)
//		return nil, ErrUsernameNotExist;
//	}
//
//	if user.Password != password {
//		beego.Debug("user password = " + user.Password + ", your password = " + password)
//		return nil, ErrPasswordNotMatched
//	}
//
//	key := &models.AdminKey{Admin:Admin{UserId:user.Id}}
//	admins, err := FindAdminVoByKey(key, false, true)
//	if err != nil {
//		beego.Error(err)
//		return nil, ErrQuery
//	}
//	if len(admins) == 0 {
//		beego.Error(ErrNotSysAdmin)
//		return nil, ErrNotSysAdmin
//	}
//	return &admins[0], nil
//}
//
//func FindAdminVoByKey(key *models.AdminKey, selectUser, selectRole bool) ([]models.AdminVo, error) {
//	admins, err := FindAdminByKey(key)
//	if err != nil {
//		beego.Error(err)
//		return []models.AdminVo{}, ErrQuery
//	}
//
//	var result []models.AdminVo
//	for _, admin := range admins {
//		adminVo := models.AdminVo{Admin:admin}
//
//		if selectUser {
//			userKey := &UserKey{User:User{Id:admin.UserId}}
//			users, err := FindUserByKey(userKey)
//			if err != nil {
//				beego.Error(err)
//			} else {
//				adminVo.User = users[0]
//			}
//		}
//
//		if selectRole {
//			roles, err := FindRoleByAdminId(admin.Id)
//			if err != nil {
//				beego.Error(err)
//			} else {
//				adminVo.Roles = roles
//			}
//		}
//
//		result = append(result, adminVo)
//	}
//	return result, nil
//}
//
//func PageAdminVo(key *models.AdminKey, selectUser bool) (*pager.Pager, error) {
//	p, err := PageAdmin(key)
//	if err != nil {
//		var admins []models.AdminVo
//		return pager.New(key.Key, 0, admins), service.ErrQuery
//	}
//
//	var adminVos []models.AdminVo
//	if admins, ok := p.PageList.([]models.Admin); ok {
//		for _, admin := range admins {
//			adminVo := models.AdminVo{Admin:admin}
//			if selectUser {
//				user, err := GetUserById(admin.UserId)
//				if err != nil {
//					beego.Error(err)
//				} else {
//					adminVo.User = *user
//				}
//			}
//			//append(adminVos, adminVo)
//		}
//	}
//	return pager.New(key.Key, p.Total, adminVos), nil
//}
//
func InsertAdminVo(admin *models.AdminVo) error {
	exist, err := userser.IsExistOfUser(&models.User{Username:admin.User.Username})
	if err != nil {
		beego.Error(err)
		return err
	} else if (exist) {
		beego.Error(ErrUsernameExist)
		return ErrUsernameExist
	}

	o := orm.NewOrm()
	o.Begin()

	now := time.Now()
	admin.User.CreateTime = now
	admin.User.UpdateTime = now
	id, err := o.Insert(admin.User)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return service.ErrInsert
	}
	admin.UserId = id
	beego.Debug(fmt.Sprintf("Add User success. userId = %v", id))

	admin.Admin.CreateTime = now
	admin.Admin.UpdateTime = now
	id, err = o.Insert(admin.Admin)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return service.ErrInsert
	}
	beego.Debug(fmt.Sprintf("Add Admin success. adminId = %v", id))
	o.Commit()
	return nil
}

func UpdateAdminVo(adminVo *models.AdminVo) error {
	if adminVo == nil {
		return service.ErrArgument
	}

	user := &models.User{
		Id:adminVo.UserId,
		Username:adminVo.Username,
	}
	exist, err := userser.IsExistOfUser(user)
	if  err != nil {
		beego.Error(err)
		return err
	} else if exist {
		beego.Error(fmt.Sprintf("username exist: username = %s", adminVo.Username))
		return ErrUsernameExist
	}

	o := orm.NewOrm();
	o.Begin()

	now := time.Now()
	adminVo.User.UpdateTime = now
	affected, err := o.Update(adminVo.User)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return service.ErrUpdate
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	adminVo.Admin.UpdateTime = now
	affected, err = o.Update(adminVo.Admin)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return service.ErrUpdate
	}
	o.Commit()
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	return nil
}

func GetAdminVoById(id int64) (*models.AdminVo, error) {
	admin, err := GetAdminById(id)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	user, err := userser.GetUserById(admin.UserId)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	return &models.AdminVo{
		Admin : admin,
		User : user,
	}, nil
}
