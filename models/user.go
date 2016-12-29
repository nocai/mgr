package models

import (
	"github.com/astaxie/beego/orm"
	"errors"
	"github.com/astaxie/beego"
	"fmt"
)

func InsertUser(user *User) error {
	if existOfUsername(user.Username) {
		beego.Debug(fmt.Sprintf("用户名存在：%v", user.Username))
		return errors.New("用户名存在")
	}

	ormer := orm.NewOrm()
	_, err := ormer.Insert(user)
	return err
}

// By Id
func GetUserById(id int64) (*User, error) {
	ormer := orm.NewOrm()

	user := &User{Id:id}
	err := ormer.Read(user)
	if err != nil {
		return nil, errors.New("查询失败")
	}
	return user, nil
}

func existOfUsername(username string) bool {
	_, err := GetUserByUsername(username)
	if err == nil {
		return true
	}
	return false
}

func GetUserByUsername(username string) (*User, error) {
	ormer := orm.NewOrm()

	user := &User{Username:username}
	err := ormer.Read(user, "Username")
	if err != nil {
		beego.Debug(err)
		return nil, ErrQuery
	}
	return user, nil;
}

func DeleteUserById(id int64) error {
	if id == 0 {
		beego.Debug("id = %v", id)
		return ErrArgument
	}

	affected, err := orm.NewOrm().Delete(&User{Id:id})
	if err != nil {
		beego.Error(err)
		return ErrDelete
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}