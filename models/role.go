package models

import (
	"time"
	"mgr/util/sqler"
	"mgr/util/key"
)

// 系统角色
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
		sqler.AppendSql(" and tmr.role_name like ?")
		sqler.AppendArg("%" + roleName + "%")
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

