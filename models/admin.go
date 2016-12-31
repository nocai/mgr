package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/util"
	"fmt"
	"time"
	"sync"
)

var (
	ErrUsernameNotExist = errors.New("用户名不存在")
	ErrUsernameExist = errors.New("用户名存在")
	ErrPasswordNotMatched = errors.New("密码错误")
	ErrNotSysAdmin = errors.New("对不起,您还不是系统管理员")
	ErrAdminNotExist = errors.New("系统管理员不存在")
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
	wg := &sync.WaitGroup{}
	if selectRole {
		wg.Add(3)
	} else {
		wg.Add(2)
	}

	admin := &Admin{}
	go func() {
		defer wg.Done()

		temp := &Admin{UserId:userId}
		err := orm.NewOrm().Read(temp, "UserId")
		if err != nil {
			beego.Error(err)
		}

		admin.Id = temp.Id
		admin.UserId = temp.UserId
		admin.AdminName = temp.AdminName
		admin.CreateTime = temp.CreateTime
		admin.UpdateTime = temp.UpdateTime
	}()
	go func() {
		defer wg.Done()
		user, err := GetUserById(userId)
		if err != nil {
			beego.Error(err)
		}
		admin.User = *user
	}()
	if selectRole {
		go func() {
			defer wg.Done()
			roles, err := FindRolesByUserId(userId)
			if err != nil {
				beego.Error(err)
			}
			admin.Roles = *roles
		}()
	}
	wg.Wait()
	beego.Debug(admin)

	return admin, nil
}

func PageAdmin(key *util.PagerKey) (*util.Pager, error) {
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
		return util.NewPager(key, 0, admins), ErrQuery
	}

	affected, err := o.Raw(sqler.GetDataSql() + key.GetOrderBySql() + key.GetLimitSql(), sqler.GetArgs()).QueryRows(&admins)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key, 0, admins), ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	ch := make(chan error, len(admins))
	go func() {
		defer close(ch)
		for _, admin := range admins {
			user, err := GetUserById(admin.UserId)
			if err != nil {
				beego.Error(err)
				ch <- err
				return
			}
			admin.User = *user
			ch <- nil
		}
	}()
	for err := range ch {
		if err != nil {
			return util.NewPager(key, 0, admins), ErrQuery
		}
	}

	return util.NewPager(key, total, admins), nil
}

func InsertAdmin(admin *Admin) error {
	if admin == nil {
		return ErrArgument
	}

	if existOfUsername(&admin.User) {
		return ErrUsernameExist
	}

	o := orm.NewOrm()
	o.Begin()

	admin.User.CreateTime = time.Now()
	admin.User.UpdateTime = time.Now()
	id, err := o.Insert(&admin.User)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return ErrInsert
	}
	admin.UserId = id

	admin.CreateTime = time.Now()
	admin.UpdateTime = time.Now()
	id, err = o.Insert(admin)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return ErrInsert
	}

	o.Commit()

	beego.Debug(fmt.Sprintf("id = %v", id))
	return nil
}

func UpdateAdmin(admin *Admin) error {
	if admin == nil {
		return ErrArgument
	}

	admin.UpdateTime = time.Now()
	admin.User.UpdateTime = time.Now()

	user := admin.User
	if existOfUsername(&user) {
		return ErrUsernameExist
	}

	o := orm.NewOrm();
	o.Begin()

	user.UpdateTime = time.Now()
	affected, err := o.Update(&user)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return ErrUpdate
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	affected, err = o.Update(admin)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return ErrUpdate
	}
	o.Commit()
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

	user, err := GetUserById(admin.UserId)
	if err != nil {
		beego.Error(err)
	}
	admin.User = *user
	return admin, nil
}

func DeleteAdminById(id int64) error {
	if id == 0 {
		return ErrArgument
	}

	admin, err := GetAdminById(id)
	if err != nil {
		return ErrAdminNotExist
	}

	o := orm.NewOrm()
	o.Begin()

	affected, err := o.Delete(&Admin{Id:id})
	if err != nil {
		beego.Error(err)
		return ErrDelete
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	affected, err = o.Delete(&User{Id:admin.UserId})
	if err != nil {
		beego.Error(err)
		return ErrDelete
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))

	o.Commit()
	return nil
}