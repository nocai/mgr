package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"mgr/util/key"
	"mgr/util/sqler"
	"strings"
	"time"
	"mgr/models/service/resser"
	"mgr/models/service/userser"
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
	orm.RegisterModelWithPrefix("t_mgr_", new(Role), new(userser.User), new(Admin), new(AdminRoleRef), new(resser.Res))

	// create table
	orm.RunSyncdb("default", false, true)

	orm.Debug = true
}
type Admin struct {
	Id        int64  `json:"id"`
	AdminName string `json:"admin_name"`
	UserId    int64  `json:"user_id"`

	CreateTime time.Time `json:"create_time"`
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
	Id       int64  `json:"id"`
	RoleName string `json:"role_name"`

	CreateTime time.Time `json:"create_time"`
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
	sqler.SetAlias("tmr")
	if this.Role != nil {
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

type AdminRoleRef struct {
	Id         int64     `json:"id"`
	AdminId    int64     `json:"admin_id"`
	RoleId     int64     `json:"role_id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (ref *AdminRoleRef) TableIndex() [][]string {
	return [][]string{
		[]string{"AdminId"},
		[]string{"RoleId"},
	}
}

type AdminRoleRefKey struct {
	*key.Key
	*AdminRoleRef

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
}

func (this *AdminRoleRefKey) NewSqler() *sqler.Sqler {
	sqler := sqler.New(this.Key)
	sqler.AppendSql("select * from t_mgr_admin_role_ref as t where 1 = 1")
	sqler.SetAlias("t")
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and t.id = ?")
		sqler.AppendArg(id)
	}
	if adminId := this.AdminId; adminId != 0 {
		sqler.AppendSql(" and t.admin_id = ?")
		sqler.AppendArg(adminId)
	}
	if roleId := this.RoleId; roleId != 0 {
		sqler.AppendSql(" and t.role_id = ?")
		sqler.AppendArg(roleId)
	}
	if !this.CreateTimeStart.IsZero() {
		sqler.AppendSql(" and t.create_time >= ?")
		sqler.AppendArg(this.CreateTimeStart)
	}
	if !this.CreateTimeEnd.IsZero() {
		sqler.AppendSql(" and t.create_time <= ?")
		sqler.AppendArg(this.CreateTimeEnd)
	}
	if !this.UpdateTimeStart.IsZero() {
		sqler.AppendSql(" and t.update_time >= ?")
		sqler.AppendArg(this.UpdateTimeStart)
	}
	if !this.UpdateTimeEnd.IsZero() {
		sqler.AppendSql(" and t.update_time <= ?")
		sqler.AppendArg(this.UpdateTimeEnd)
	}

	return sqler
}
