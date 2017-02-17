package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"mgr/util"
)

type User struct {
	ModelBase

	Id       int64
	Username string `orm:"unique"`
	Password string
}

// 多字段索引
func (user *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Username"},
	}
}

// 多字段唯一键
func (user *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Username"},
	}
}

func InsertUser(user *User) error {
	if user == nil {
		beego.Error("user is nil")
		return ErrArgument
	}
	return insertUser(nil, user)
}

func insertUser(o orm.Ormer, user *User) error {
	if o == nil {
		o = orm.NewOrm()
	}

	if err := checkUser(user); err != nil {
		return err
	}
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()

	ormer := orm.NewOrm()
	id, err := ormer.Insert(user)
	if err != nil {
		beego.Error(err)
		return ErrInsert
	}
	beego.Debug(fmt.Sprintf("id = %v", id))
	return nil

}

func checkUser(user *User) error {
	if user == nil {
		beego.Error("user is nil")
		return ErrArgument
	}

	if user.Username == "" {
		beego.Error("user.Username is nil")
		return ErrArgument
	}
	if existOfUsername(user) {
		return ErrArgument
	}

	if user.Password == "" {
		beego.Error("user.Password is nil")
		return ErrArgument
	}
	return nil
}

// By Id
func GetUserById(id int64) (*User, error) {
	if id == 0 {
		beego.Error("id = 0")
		return nil, ErrArgument
	}

	ormer := orm.NewOrm()

	user := &User{Id:id}
	err := ormer.Read(user)
	if err != nil {
		return nil, ErrQuery
	}
	return user, nil
}

func existOfUsername(user *User) bool {
	if temp, err := GetUserByUsername(user.Username); err == nil {
		if temp.Id != user.Id {
			beego.Info(fmt.Sprintf("用户名存在%v", user.Username))
			return true
		}
	}
	return false
}

func GetUserByUsername(username string) (*User, error) {
	user := User{Username:username}
	key := &UserKey{User:user}
	users, err := ListUserByKey(key)
	if err != nil {
		if len(users) > 1 {
			beego.Error(fmt.Sprintf("useranme重复：username = %v, 重复数据 = %v", username, users))
		}
		return &users[0], nil
	}

	return &user, nil;
}

type UserKey struct {
	util.PagerKey
	User

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
}

func (key *UserKey) generateSql() {
	if key.IsEmptySql() {
		key.AppendSql(`select * from t_mgr_user as tmu where 1 = 1`)

		if id := key.Id; id != 0 {
			key.AppendSql(" and tmu.id = ?")
			key.AppendArg(id)
		}
		if username := key.Username; username != "" {
			key.AppendSql(" and tmu.username = ?")
			key.AppendArg(key.Username)
		}
		if password := key.Password; password != "" {
			key.AppendSql(" and tmu.password = ?")
			key.AppendArg(password)
		}
		if createTime := key.CreateTime; !createTime.IsZero() {
			key.AppendSql(" and tmu.create_time = ?")
			key.AppendArg(createTime)
		}
		if updateTime := key.UpdateTime; !updateTime.IsZero() {
			key.AppendSql(" and tmu.update_time = ?")
			key.AppendArg(updateTime)
		}

		if createTimeStart := key.CreateTimeStart; !createTimeStart.IsZero() {
			key.AppendSql(" and tmu.create_time >= ?")
			key.AppendArg(createTimeStart)
		}
		if createTimeEnd := key.CreateTimeEnd; !createTimeEnd.IsZero() {
			key.AppendSql(" and tmu.create_time < ?")
			key.AppendArg(createTimeEnd)
		}

		if updateTimeStart := key.UpdateTimeStart; !updateTimeStart.IsZero() {
			key.AppendSql(" and tmu.update_time >= ?")
			key.AppendArg(updateTimeStart)
		}
		if updateTimeEnd := key.UpdateTimeEnd; !updateTimeEnd.IsZero() {
			key.AppendSql(" and tmu.update_time < ?")
			key.AppendArg(updateTimeEnd)
		}
		if keyWord := key.KeyWord; keyWord != "" {
			key.AppendSql(" and tmu.username like ?")
			key.AppendArg("%" + key.KeyWord + "%")
		}

		key.AppendSql(" and tmu.invalid = ?")
		key.AppendArg(key.Invalid)
	}
}

func PageUser(key *UserKey) (*util.Pager, error) {
	total, err := CountUserByKey(key)
	if err != nil {
		return nil, err
	}
	users, err := ListUserByKey(key)
	if err != nil {
		return nil, err
	}
	return util.NewPager(&key.PagerKey, total, users), nil
}

func CountUserByKey(key *UserKey) (int64, error) {
	if key.IsEmptySql() {
		key.generateSql()
	}

	o := orm.NewOrm()

	var total int64
	err := o.Raw(key.GetCountSql(), key.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return 0, ErrQuery
	}
	return total, nil
}

func ListUserByKey(key *UserKey) ([]User, error) {
	beego.Error(fmt.Sprintf("%+v", key))
	if key.IsEmptySql() {
		key.generateSql()
	}

	o := orm.NewOrm()

	var users []User
	affected, err := o.Raw(key.GetSql(), key.GetArgs()).QueryRows(&users)
	if err != nil {
		beego.Error(err)
		return users, ErrQuery
	}
	beego.Debug("affected = %v", affected)
	return users, nil

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