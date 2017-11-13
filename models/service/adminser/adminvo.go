package adminser

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"mgr/models"
	"mgr/models/service"
	"mgr/models/service/userser"
	"mgr/util/key"
	"mgr/util/pager"
	"mgr/util/sqler"
	"strings"
	"time"
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

type AdminVo struct {
	*models.Admin
	*models.User `json:"user"`

	Roles []models.Role `json:"roles"`
}

type AdminVoKey struct {
	*key.Key
	*models.Admin
	Invalid models.ValidEnum
}

func (this *AdminVoKey) NewSqler() *sqler.Sqler {
	sqler := sqler.New(this.Key)

	sqler.AppendSql(`select tma.* from t_mgr_admin as tma join t_mgr_user as tmu on tma.user_id = tmu.id where 1 = 1`)
	if this.Invalid != models.ValidAll {
		sqler.AppendSql(" and tmu.invalid = ?")
		sqler.AppendArg(this.Invalid)
	}
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tma.id = ?")
		sqler.AppendArg(id)
	}
	if adminName := this.AdminName; adminName != "" {
		sqler.AppendSql(" and tma.admin_name")
		if strings.Contains(adminName, "%") {
			sqler.AppendSql(" like ?")
		} else {
			sqler.AppendSql(" = ?")
		}
		sqler.AppendArg(adminName)
	}
	if userId := this.UserId; userId != 0 {
		sqler.AppendSql(" and tma.user_id = ?")
		sqler.AppendArg(userId)
	}

	//if createTimeStart := this.CreateTimeStart; !createTimeStart.IsZero() {
	//	sqler.AppendSql(" and tma.create_time >= ?")
	//	sqler.AppendArg(createTimeStart)
	//}
	//if createTimeEnd := this.CreateTimeEnd; !createTimeEnd.IsZero() {
	//	sqler.AppendSql(" and tma.create_time < ?")
	//	sqler.AppendArg(createTimeEnd)
	//}
	//
	//if updateTimeStart := this.UpdateTimeStart; !updateTimeStart.IsZero() {
	//	sqler.AppendSql(" and tma.update_time >= ?")
	//	sqler.AppendArg(updateTimeStart)
	//}
	//if updateTimeEnd := this.UpdateTimeEnd; !updateTimeEnd.IsZero() {
	//	sqler.AppendSql(" and tma.update_time < ?")
	//	sqler.AppendArg(updateTimeEnd)
	//}
	//if keyWord := this.KeyWord; keyWord != "" {
	//	sqler.AppendSql(" and tma.admin_name like ?")
	//	sqler.AppendArg("%" + keyWord + "%")
	//}

	return sqler
}

func PageAdminVo(key *AdminVoKey) (*pager.Pager, error) {
	sqler := key.NewSqler()
	o := orm.NewOrm()

	var total int64
	err := o.Raw(sqler.GetCountSqlAndArgs()).QueryRow(&total)
	if err != nil {
		return pager.New(key.Key, 0, []models.Admin{}), errors.Wrap(err, service.MsgQuery)
	}

	var admins []models.Admin
	affected, err := o.Raw(sqler.GetSqlAndArgs()).QueryRows(&admins)
	if err != nil {
		beego.Error(err)
		return pager.New(key.Key, 0, []models.Admin{}), errors.Wrap(err, service.MsgQuery)
	}
	beego.Debug(fmt.Sprintf("affected = %d", affected))
	if affected == 0 {
		return pager.New(key.Key, 0, []models.Admin{}), nil
	}

	var adminVos []AdminVo
	for i, _ := range admins {
		adminVos = append(adminVos, AdminVo{Admin: &admins[i]})
	}

	ch := make(chan error, len(adminVos))
	for index, _ := range adminVos {
		go func(i int, ch chan error) {
			user, err := userser.GetUserById(adminVos[i].UserId)
			if err != nil {
				ch <- err
			} else {
				adminVos[i].User = user
				ch <- nil
			}
		}(index, ch)
	}

	for i := 0; i < len(adminVos); i++ {
		err := <-ch
		if err != nil {
			return pager.New(key.Key, 0, []models.Admin{}), err
		}
	}
	return pager.New(key.Key, total, adminVos), nil
}

func InsertAdminVo(admin *AdminVo) error {
	exist, err := userser.IsExistOfUser(&models.User{Username: admin.User.Username})
	if err != nil {
		return errors.Wrap(err, service.MsgInsert)
	} else if exist {
		return ErrUsernameExist
	}

	o := orm.NewOrm()
	o.Begin()

	now := time.Now()
	admin.User.CreateTime = now
	admin.User.UpdateTime = now
	id, err := o.Insert(admin.User)
	if err != nil {
		o.Rollback()
		return errors.Wrap(err, service.MsgInsert)
	}
	admin.UserId = id
	beego.Debug(fmt.Sprintf("Add User success. userId = %v", id))

	admin.Admin.CreateTime = now
	admin.Admin.UpdateTime = now
	id, err = o.Insert(admin.Admin)
	if err != nil {
		o.Rollback()
		return errors.Wrap(err, service.MsgInsert)
	}
	beego.Debug(fmt.Sprintf("Add Admin success. adminId = %v", id))
	o.Commit()
	return nil
}

func UpdateAdminVo(adminVo *AdminVo) error {
	if adminVo == nil {
		return service.ErrArgument
	}

	user := &models.User{
		Id:       adminVo.UserId,
		Username: adminVo.Username,
	}
	exist, err := userser.IsExistOfUser(user)
	if err != nil {
		beego.Error(err)
		return err
	} else if exist {
		beego.Error(fmt.Sprintf("username exist: username = %s", adminVo.Username))
		return ErrUsernameExist
	}

	o := orm.NewOrm()
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

func GetAdminVoById(id int64) (*AdminVo, error) {
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

	return &AdminVo{
		Admin: admin,
		User:  user,
	}, nil
}
