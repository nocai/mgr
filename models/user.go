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
	*util.Key
	User

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *UserKey) getSqler() *util.Sqler {
	sqler := &util.Sqler{Key:this.Key}
	sqler.AppendSql(`select * from t_mgr_user as tmu where 1 = 1`)

	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tmu.id = ?")
		sqler.AppendArg(id)
	}
	if username := this.Username; username != "" {
		sqler.AppendSql(" and tmu.username = ?")
		sqler.AppendArg(this.Username)
	}
	if password := this.Password; password != "" {
		sqler.AppendSql(" and tmu.password = ?")
		sqler.AppendArg(password)
	}
	if createTime := this.CreateTime; !createTime.IsZero() {
		sqler.AppendSql(" and tmu.create_time = ?")
		sqler.AppendArg(createTime)
	}
	if updateTime := this.UpdateTime; !updateTime.IsZero() {
		sqler.AppendSql(" and tmu.update_time = ?")
		sqler.AppendArg(updateTime)
	}

	if createTimeStart := this.CreateTimeStart; !createTimeStart.IsZero() {
		sqler.AppendSql(" and tmu.create_time >= ?")
		sqler.AppendArg(createTimeStart)
	}
	if createTimeEnd := this.CreateTimeEnd; !createTimeEnd.IsZero() {
		sqler.AppendSql(" and tmu.create_time < ?")
		sqler.AppendArg(createTimeEnd)
	}

	if updateTimeStart := this.UpdateTimeStart; !updateTimeStart.IsZero() {
		sqler.AppendSql(" and tmu.update_time >= ?")
		sqler.AppendArg(updateTimeStart)
	}
	if updateTimeEnd := this.UpdateTimeEnd; !updateTimeEnd.IsZero() {
		sqler.AppendSql(" and tmu.update_time < ?")
		sqler.AppendArg(updateTimeEnd)
	}
	if keyWord := this.KeyWord; keyWord != "" {
		sqler.AppendSql(" and tmu.username like ?")
		sqler.AppendArg("%" + this.KeyWord + "%")
	}

	sqler.AppendSql(" and tmu.invalid = ?")
	sqler.AppendArg(this.Invalid)
	return sqler;
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
	return util.NewPager(key.Key, total, users), nil
}

func CountUserByKey(key *UserKey) (int64, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var total int64
	err := o.Raw(sqler.GetCountSql(), sqler.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return 0, ErrQuery
	}
	return total, nil
}

func ListUserByKey(key *UserKey) ([]User, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var users []User
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&users)
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