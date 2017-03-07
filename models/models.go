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
	ErrDataDuplication = errors.New("数据重复")
)

// 所有模型共同的属性
type ModelBase struct {
	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func init() {
	orm.RegisterDataBase("default", "mysql", "saas20:sass20@tcp(columbus.unovo.com.cn:3306)/mgr?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	// register model
	//orm.RegisterModel(new(User), new(Role), new(AdminRoleRef), new(Admin))

	orm.RegisterModelWithPrefix("t_mgr_", new(User), new(Role), new(AdminRoleRef), new(Admin), new(Res))
	// create table
	orm.RunSyncdb("default", false, true)

	orm.Debug = true
}
