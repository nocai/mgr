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

type Admin struct {
	ModelBase

	Id        int64 `json:"id"`
	AdminName string `json:"admin_name"`
	UserId    int64 `orm:"unique" json:"user_id"`
}

type AdminVo struct {
	Admin

	User  User `orm:"-" json:"user"`
	Roles []Role `orm:"-" json:"roles"`
}

func Login(username, password string) (*AdminVo, error) {
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

func GetAdminByUserId(userId int64, selectRole bool) (*AdminVo, error) {
	wg := &sync.WaitGroup{}
	if selectRole {
		wg.Add(3)
	} else {
		wg.Add(2)
	}

	admin := &AdminVo{}
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

type AdminKey struct {
	util.PagerKey
	Admin

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
}

func (this *AdminKey) generateSql() {
	if this.IsEmptySql() {
		this.AppendSql(`select * from t_mgr_admin as tma where 1 = 1`)

		if id := this.Id; id != 0 {
			this.AppendSql(" and tma.id = ?")
			this.AppendArg(id)
		}
		if adminName := this.AdminName; adminName != "" {
			this.AppendSql(" and tma.admin_name = ?")
			this.AppendArg(adminName)
		}
		if userId := this.UserId; userId != 0 {
			this.AppendSql(" and tma.user_id = ?")
			this.AppendArg(userId)
		}

		if createTimeStart := this.CreateTimeStart; !createTimeStart.IsZero() {
			this.AppendSql(" and tma.create_time >= ?")
			this.AppendArg(createTimeStart)
		}
		if createTimeEnd := this.CreateTimeEnd; !createTimeEnd.IsZero() {
			this.AppendSql(" and tma.create_time < ?")
			this.AppendArg(createTimeEnd)
		}

		if updateTimeStart := this.UpdateTimeStart; !updateTimeStart.IsZero() {
			this.AppendSql(" and tma.update_time >= ?")
			this.AppendArg(updateTimeStart)
		}
		if updateTimeEnd := this.UpdateTimeEnd; !updateTimeEnd.IsZero() {
			this.AppendSql(" and tma.update_time < ?")
			this.AppendArg(updateTimeEnd)
		}
		if keyWord := this.KeyWord; keyWord != "" {
			this.AppendSql(" and tma.admin_name like ?")
			this.AppendArg("%" + keyWord + "%")
		}

		this.AppendSql(" and tma.invalid = ?")
		this.AppendArg(this.Invalid)
	}
}

func CountAdminByKey(key *AdminKey) (int64, error) {
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

func ListAdminByKey(key *AdminKey) ([]Admin, error) {
	if key.IsEmptySql() {
		key.generateSql()
	}

	o := orm.NewOrm()

	var admins []Admin
	affected, err := o.Raw(key.GetSql(), key.GetArgs()).QueryRows(&admins)
	if err != nil {
		beego.Error(err)
		return admins, err
	}
	beego.Debug(fmt.Sprintf("affected = %v, %+v", affected, admins))
	return admins, nil
}

func PageAdmin(key *AdminKey) (*util.Pager, error) {
	total, err := CountAdminByKey(key)
	if err != nil {
		var admins []Admin
		return util.NewPager(&key.PagerKey, 0, admins), ErrQuery
	}

	admins, err := ListAdminByKey(key)
	if err != nil {
		return util.NewPager(&key.PagerKey, 0, admins), ErrQuery
	}

	return util.NewPager(&key.PagerKey, total, admins), nil
}

func PageAdminVo(key *AdminKey, selectUser bool) (*util.Pager, error) {
	pager, err := PageAdmin(key)
	if err != nil {
		var admins []AdminVo
		return util.NewPager(&key.PagerKey, 0, admins), ErrQuery
	}

	var adminVos []AdminVo
	if admins, ok := pager.PageList.([]Admin); ok {
		for _, admin := range admins {
			adminVo := AdminVo{Admin:admin}
			if selectUser {
				user, err := GetUserById(admin.UserId)
				if err != nil {
					beego.Error(err)
				} else {
					adminVo.User = *user
				}
			}
			//append(adminVos, adminVo)
		}
	}
	return util.NewPager(&key.PagerKey, pager.Total, adminVos), nil
}

func InsertAdmin(admin *Admin) error {
	return insertAdmin(nil, admin)
}

func insertAdmin(o orm.Ormer, admin *Admin) error {
	if o == nil {
		o = orm.NewOrm()
	}

	if admin.UserId == 0 {
		beego.Error("管理员的UserId必须填")
		return ErrArgument
	}

	now := time.Now()
	if admin.CreateTime.IsZero() {
		admin.CreateTime = now
	}
	if admin.UpdateTime.IsZero() {
		admin.UpdateTime = now
	}

	id, err := o.Insert(admin)
	if err != nil {
		beego.Error(err)
		return ErrInsert
	}
	beego.Debug(fmt.Sprintf("添加Admin, id = %v", id))
	return nil
}

func InsertAdminVo(admin *AdminVo) error {
	if admin == nil {
		return ErrArgument
	}

	if existOfUsername(&admin.User) {
		return ErrUsernameExist
	}

	o := orm.NewOrm()
	o.Begin()

	now := time.Now()

	admin.User.CreateTime = now
	admin.User.UpdateTime = now
	id, err := o.Insert(&admin.User)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return ErrInsert
	}
	admin.UserId = id
	beego.Debug(fmt.Sprintf("UserId = %v", id))

	admin.CreateTime = now
	admin.UpdateTime = now
	err = insertAdmin(o, &admin.Admin)
	if err != nil {
		o.Rollback()
		return ErrInsert
	}

	o.Commit()

	return nil
}

func UpdateAdmin(admin *AdminVo) error {
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

func GetAdminById(id int64) (*AdminVo, error) {
	if id == 0 {
		return nil, ErrArgument
	}

	admin := &AdminVo{Admin:Admin{Id:id}}
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