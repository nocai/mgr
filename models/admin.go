package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/util"
	"fmt"
	"time"
)

var (
	ErrUsernameNotExist = errors.New("用户名不存在")
	ErrPasswordNotMatched = errors.New("密码错误")
	ErrNotSysAdmin = errors.New("对不起,您还不是系统管理员")
)

func Login(username, password string) (*Admin, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		beego.Debug(err)
		return nil, ErrUsernameNotExist;
	}

	if user.Password != password {
		beego.Debug("user password = " + user.Password + ", your password = " + password)
		return nil, ErrPasswordNotMatched
	}

	admin, err := GetAdminByUserId(user.Id, false)
	if err != nil {
		beego.Debug(err)
		return admin, ErrNotSysAdmin
	}
	return admin, nil
}

func GetAdminByUserId(userId int64, selectRole bool) (*Admin, error) {
	ormer := orm.NewOrm()

	admin := &Admin{UserId:userId}
	err := ormer.Read(admin, "UserId")
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	if selectRole {
		roles, err := FindRolesByUserId(userId)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		admin.Roles = *roles
	}

	return admin, nil
}

func PageAdmin(key *util.PagerKey) *util.Pager {
	sqler := util.NewSqler(`select tma.* from t_mgr_admin as tma where 1 = 1`)

	if adminName, ok := key.Data["adminName"].(string); ok && adminName != "" {
		sqler.AppendDataSql(" and tma.admin_name like ?")
		sqler.AppendArg("%" + adminName + "%")
	}

	o := orm.NewOrm()

	var total int64
	var admins []Admin
	err := o.Raw(sqler.GetCountSql(), sqler.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key, 0, admins)
	}

	affected, err := o.Raw(sqler.GetDataSql(), sqler.GetArgs()).QueryRows(&admins)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key, 0, admins)
	}

	beego.Info(fmt.Sprintf("affected = %v", affected))
	return util.NewPager(key, total, admins)
}

func AddAdmin(admin *Admin) error {
	if admin == nil {
		return ErrArgument
	}

	admin.User.CreateTime = time.Now()
	admin.User.UpdateTime = time.Now()
	err := InsertUser(&admin.User)
	if err != nil {
		return err
	}
	admin.UserId = admin.User.Id

	admin.CreateTime = time.Now()
	admin.UpdateTime = time.Now()

	o := orm.NewOrm()
	affected, err := o.Insert(admin)
	if err != nil {
		beego.Error(err)
		return ErrInsert
	}


	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}

func GetAdminById(id int64) (*Admin, error) {
	if id == 0 {
		return nil, ErrArgument
	}

	admin := &Admin{Id:id}
	err := orm.NewOrm().Read(admin)
	if err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}

	return admin, nil
}

func DeleteAdminById(id int64) error {
	if id == 0 {
		return ErrArgument
	}

	admin, err := GetAdminById(id)
	if err != nil {
		return ErrDelete
	}
	err = DeleteUserById(admin.UserId)
	if err != nil {
		return ErrDelete
	}

	o := orm.NewOrm()
	affected, err := o.Delete(admin)
	if err != nil {
		beego.Error(err)
		return ErrDelete
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}