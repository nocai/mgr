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

	key := &AdminKey{Admin:Admin{UserId:user.Id}}
	admins, err := FindAdminVoByKey(key, false, true)
	if err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}
	if len(admins) == 0 {
		beego.Error(ErrNotSysAdmin)
		return nil, ErrNotSysAdmin
	}
	return &admins[0], nil
}

func FindAdminVoByKey(key *AdminKey, selectUser, selectRole bool) ([]AdminVo, error) {
	admins, err := FindAdminByKey(key)
	if err != nil {
		beego.Error(err)
		return []AdminVo{}, ErrQuery
	}

	var result []AdminVo
	for _, admin := range admins {
		adminVo := AdminVo{Admin:admin}

		if selectUser {
			userKey := &UserKey{User:User{Id:admin.UserId}}
			users, err := FindUserByKey(userKey)
			if err != nil {
				beego.Error(err)
			} else {
				adminVo.User = users[0]
			}
		}

		if selectRole {
			roles, err := FindRoleByAdminId(admin.Id)
			if err != nil {
				beego.Error(err)
			} else {
				adminVo.Roles = roles
			}
		}

		result = append(result, adminVo)
	}
	return result, nil
}

type AdminKey struct {
	*util.Key
	Admin

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *AdminKey) getSqler() *util.Sqler {
	sqler := &util.Sqler{Key:this.Key}

	sqler.AppendSql(`select * from t_mgr_admin as tma where 1 = 1`)
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tma.id = ?")
		sqler.AppendArg(id)
	}
	if adminName := this.AdminName; adminName != "" {
		sqler.AppendSql(" and tma.admin_name = ?")
		sqler.AppendArg(adminName)
	}
	if userId := this.UserId; userId != 0 {
		sqler.AppendSql(" and tma.user_id = ?")
		sqler.AppendArg(userId)
	}

	if createTimeStart := this.CreateTimeStart; !createTimeStart.IsZero() {
		sqler.AppendSql(" and tma.create_time >= ?")
		sqler.AppendArg(createTimeStart)
	}
	if createTimeEnd := this.CreateTimeEnd; !createTimeEnd.IsZero() {
		sqler.AppendSql(" and tma.create_time < ?")
		sqler.AppendArg(createTimeEnd)
	}

	if updateTimeStart := this.UpdateTimeStart; !updateTimeStart.IsZero() {
		sqler.AppendSql(" and tma.update_time >= ?")
		sqler.AppendArg(updateTimeStart)
	}
	if updateTimeEnd := this.UpdateTimeEnd; !updateTimeEnd.IsZero() {
		sqler.AppendSql(" and tma.update_time < ?")
		sqler.AppendArg(updateTimeEnd)
	}
	if keyWord := this.KeyWord; keyWord != "" {
		sqler.AppendSql(" and tma.admin_name like ?")
		sqler.AppendArg("%" + keyWord + "%")
	}

	sqler.AppendSql(" and tma.invalid = ?")
	sqler.AppendArg(this.Invalid)
	return sqler
}

func CountAdminByKey(key *AdminKey) (int64, error) {
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

func FindAdminByKey(key *AdminKey) ([]Admin, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var admins []Admin
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&admins)
	if err != nil {
		beego.Error(err)
		return admins, err
	}
	beego.Debug(fmt.Sprintf("affected = %v, %+v", affected, admins))
	if affected == 0 {
		return []Admin{}, nil
	}
	return admins, nil
}

func PageAdmin(key *AdminKey) (*util.Pager, error) {
	total, err := CountAdminByKey(key)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, []Admin{}), ErrQuery
	}

	admins, err := FindAdminByKey(key)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, []Admin{}), ErrQuery
	}
	return util.NewPager(key.Key, total, admins), nil
}

func PageAdminVo(key *AdminKey, selectUser bool) (*util.Pager, error) {
	pager, err := PageAdmin(key)
	if err != nil {
		var admins []AdminVo
		return util.NewPager(key.Key, 0, admins), ErrQuery
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
	return util.NewPager(key.Key, pager.Total, adminVos), nil
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

	if isExistOfUser(&admin.User) {
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
	if isExistOfUser(&user) {
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