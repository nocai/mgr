package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func init() {
	//orm.RegisterDataBase("defaultlt", "mysql", "saas20:sass20@tcp(columbus.unovo.com.cn:3306)/mgr?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "devel:devel@tcp(139.196.191.164:3306)/mgr?charset=utf8", 30)
	orm.RegisterDataBase("default", "mysql", "root:root@/mgr?charset=utf8", 30)
	// register model
	//orm.RegisterModel(new(User), new(Role), new(AdminRoleRef), new(Admin))

	//orm.RegisterModelWithPrefix("t_mgr_", new(User), new(role.Role), new(AdminRoleRef), new(Admin), new(Res))
	orm.RegisterModelWithPrefix("t_mgr_", new(Role))

	// create table
	orm.RunSyncdb("default", false, true)

	orm.Debug = true
}


// 所有模型共同的属性
type ModelBase struct {
	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
