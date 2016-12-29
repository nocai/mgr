package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"errors"
)

var (
	ErrQuery = errors.New("查询失败")
	ErrInsert = errors.New("添加失败")
	ErrUpdate = errors.New("更新失败")
	ErrDelete = errors.New("删除失败")

	ErrArgument = errors.New("无效参数")
)

type Admin struct {
	Id         int64 `json:"id"`
	AdminName  string `json:"admin_name"`
	UserId     int64 `orm:"unique" json:"user_id"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`

	User       User `orm:"-" json:"user"`
	Roles      []Role `orm:"-" json:"roles"`
}

func (Admin *Admin) TableName() string {
	return "t_mgr_admin"
}

type Role struct {
	Id         int64  `json:"id"`
	RoleName   string `orm:"unique" json:"role_name"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type User struct {
	Id         int64
	Username   string `orm:"unique"`
	Password   string

	CreateTime time.Time
	UpdateTime time.Time
}

type AdminRoleRef struct {
	Id      int64
	AdminId int64
	RoleId  int64
}

func (ref *AdminRoleRef) TableName() string {
	return "t_mgr_admin_role_ref"
}

func (ref *AdminRoleRef) TableIndex() [][]string {
	return [][] string{
		[] string{"AdminId"},
		[] string{"RoleId"},
	}
}

func (user *User) TableName() string {
	return "t_mgr_user"
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
func (role *Role) TableName() string {
	return "t_mgr_role"
}

// 多字段索引
func (role *Role) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "RoleName"},
	}
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:wWXdjF9r0iGgwSKY@tcp(139.196.152.74:3306)/test?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	// register model
	orm.RegisterModel(new(User), new(Role), new(AdminRoleRef), new(Admin))

	//orm.RegisterModelWithPrefix("t_mgr_", new(User), new(Role))
	// create table
	orm.RunSyncdb("default", false, true)

	orm.Debug = true
}
