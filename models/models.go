package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	//"github.com/astaxie/beego"
	"mgr/util/sqler"
	"mgr/util/key"
	"strings"
)

func init() {
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	//mysqluser := beego.AppConfig.String("mysqluser")
	//mysqlpass := beego.AppConfig.String("mysqlpass")
	//beego.Info("mysqluser:" , mysqluser)
	//mysqlurls := beego.AppConfig.String("mysqlurls")
	//mysqldb := beego.AppConfig.String("mysqldb")
	//orm.RegisterDataBase("default", "mysql", mysqluser + ":" + mysqlpass +"@/" + mysqldb + "?charset=utf8", 30)
	orm.RegisterDataBase("default", "mysql", "root:root@/mgr?charset=utf8", 30)
	// register model
	//orm.RegisterModelWithPrefix("t_mgr_", new(User), new(role.Role), new(AdminRoleRef), new(Admin), new(Res))
	orm.RegisterModelWithPrefix("t_mgr_", new(Role), new(User), new(Admin))

	// create table
	orm.RunSyncdb("default", false, true)

	orm.Debug = true
}

type ValidEnum int

const (
	// 无效的
	Invalid ValidEnum = iota
	// 有效的
	Valid
	// 所有
	ValidAll
)

type User struct {
	Id         int64
	Username   string
	Password   string

	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`

	Invalid    ValidEnum `json:"invalid"`
}


// 多字段唯一键
func (user *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Username"},
	}
}

// 多字段索引
func (user *User) TableIndex() [][]string {
	return [][]string{
		[]string{"CreateTime"},
		[]string{"UpdateTime"},
		[]string{"Invalid"},
	}
}


type UserKey struct {
	*key.Key
	*User

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *UserKey) NewSqler() *sqler.Sqler {
	sqler := sqler.New(this.Key)
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
	if this.Invalid != ValidAll {
		sqler.AppendSql(" and tmu.invalid = ?")
		sqler.AppendArg(this.Invalid)
	}

	return sqler;
}

type Admin struct {
	Id         int64 `json:"id"`
	AdminName  string `json:"admin_name"`
	UserId     int64 `json:"user_id"`

	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// 多字段唯一键
func (admin *Admin) TableUnique() [][]string {
	return [][]string{
		[]string{"AdminName"},
		[]string{"UserId"},
	}
}
// 多字段索引
func (admin *Admin) TableIndex() [][]string {
	return [][]string{
		[]string{"CreateTime"},
		[]string{"UpdateTime"},
	}
}

type AdminKey struct {
	*key.Key
	*Admin

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *AdminKey) NewSqler() *sqler.Sqler {
	sqler := sqler.New(this.Key)

	sqler.AppendSql(`select * from t_mgr_admin as tma where 1 = 1`)
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

	return sqler
}



// 系统角色
type Role struct {
	Id         int64  `json:"id"`
	RoleName   string `json:"role_name"`

	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// 多字段唯一键
func (role *Role) TableUnique() [][]string {
	return [][]string{
		[]string{"RoleName"},
	}
}
// 多字段索引
func (role *Role) TableIndex() [][]string {
	return [][]string{
		[]string{"CreateTime"},
		[]string{"UpdateTime"},
	}
}

type RoleKey struct {
	*key.Key
	*Role

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *RoleKey) NewSqler() *sqler.Sqler {
	sqler := sqler.New(this.Key)
	sqler.AppendSql("select * from t_mgr_role as tmr where 1 = 1")
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tmr.id = ?")
		sqler.AppendArg(id)
	}
	if roleName := this.RoleName; roleName != "" {
		sqler.AppendSql(" and tmr.role_name ")
		if strings.Contains(roleName, "%") {
			sqler.AppendSql(" like ?")
		} else {
			sqler.AppendSql(" = ?")
		}

		sqler.AppendArg(roleName)
	}
	if !this.CreateTimeStart.IsZero() {
		sqler.AppendSql(" and tmr.create_time >= ?")
		sqler.AppendArg(this.CreateTimeStart)
	}
	if !this.CreateTimeEnd.IsZero() {
		sqler.AppendSql(" and tmr.create_time <= ?")
		sqler.AppendArg(this.CreateTimeEnd)
	}
	if !this.UpdateTimeStart.IsZero() {
		sqler.AppendSql(" and tmr.update_time >= ?")
		sqler.AppendArg(this.UpdateTimeStart)
	}
	if !this.UpdateTimeEnd.IsZero() {
		sqler.AppendSql(" and tmr.update_time <= ?")
		sqler.AppendArg(this.UpdateTimeEnd)
	}

	return sqler
}
