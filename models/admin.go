package models

import (
	"time"
	"mgr/util/sqler"
	"mgr/util/key"
)

type Admin struct {
	ModelBase

	Id        int64 `json:"id"`
	AdminName string `json:"admin_name"`
	UserId    int64 `orm:"unique" json:"user_id"`

	Invalid   bool `json:"invalid"`
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

type AdminVo struct {
	*Admin
	*User `orm:"-" json:"user"`

	Roles []Role `orm:"-" json:"roles"`
}
