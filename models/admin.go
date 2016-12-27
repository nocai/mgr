package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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