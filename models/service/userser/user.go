package userser

import (
	"time"
	"mgr/util/sqler"
	"mgr/util/key"
)

type ValidEnum int

const (
	// 所有
	_ ValidEnum = iota
	// 无效的
	Invalid
	// 有效的
	Valid
)

type User struct {
	Id       int64
	Username string
	Password string

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`

	Invalid ValidEnum `json:"invalid"`
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
	if int(this.Invalid) == 1 || this.Invalid == 2 {
		sqler.AppendSql(" and tmu.invalid = ?")
		sqler.AppendArg(this.Invalid)
	}

	return sqler
}

