package models

import (
	"mgr/util"
	"time"
)

type Role struct {
	ModelBase

	Id       int64  `json:"id"`
	RoleName string `orm:"unique" json:"role_name"`
}

// 多字段索引
func (role *Role) TableIndex() [][]string {
	return [][]string{
		[]string{"RoleName"},
	}
}

type RoleKey struct {
	*util.Key

	Role
	CreateTimeStart time.Time
	CreateTimeEnd   time.Time
	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
	KeyWord         string
}

func (this *RoleKey) GetSqler() *util.Sqler {
	sqler := this.NewSqler()

	sqler.AppendSql("select * from t_mgr_role as tmr where 1 = 1")
	if id := this.Id; id != 0 {
		sqler.AppendSql(" and tmr.id = ?")
		sqler.AppendArg(id)
	}
	if roleName := this.RoleName; roleName != "" {
		sqler.AppendSql(" and tmr.role_name like ?")
		sqler.AppendArg("%" + roleName + "%")
	}
	return sqler
}