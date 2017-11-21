package resser

import (
	"time"
	"mgr/util/sqler"
	"mgr/util/key"
)

type ResKey struct {
	*key.Key

	Res

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time

	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
}

func (this *ResKey) getSqler() *sqler.Sqler {
	sqler := sqler.New(this.Key)

	sqler.AppendSql(`select * from t_mgr_res as tmr where 1 = 1`)
	if this.Id != 0 {
		sqler.AppendSql(" and tmr.id = ?")
		sqler.AppendArg(this.Id)
	}
	if this.ResName != "" {
		sqler.AppendSql(" and tmr.res_name = ?")
		sqler.AppendArg(this.ResName)
	}
	if this.Path != "" {
		sqler.AppendSql(" and tmr.path = ?")
		sqler.AppendArg(this.Path)
	}
	if resType := this.ResType; resType != 0 {
		sqler.AppendSql(" and tmr.res_type = ?")
		sqler.AppendArg(resType)
	}
	if pid := this.Pid; pid != 0 {
		sqler.AppendSql(" and tmr.pid = ?")
		sqler.AppendArg(pid)
	}
	if createTime := this.CreateTime; !createTime.IsZero() {
		sqler.AppendSql(" and tmr.create_time = ?")
		sqler.AppendArg(createTime)
	}
	if updateTime := this.UpdateTime; !updateTime.IsZero() {
		sqler.AppendSql(" and tmr.update_time = ?")
		sqler.AppendArg(updateTime)
	}
	if createTimeStart := this.CreateTimeStart; !createTimeStart.IsZero() {
		sqler.AppendSql(" and tmr.create_time >= ?")
		sqler.AppendArg(createTimeStart)
	}
	if createTimeEnd := this.CreateTimeEnd; !createTimeEnd.IsZero() {
		sqler.AppendSql(" and tmr.create_time <= ?")
		sqler.AppendArg(createTimeEnd)
	}
	if updateTimeStart := this.UpdateTimeStart; !updateTimeStart.IsZero() {
		sqler.AppendSql(" and tmr.update_time >= ?")
		sqler.AppendArg(updateTimeStart)
	}
	if updateTimeEnd := this.UpdateTimeStart; !updateTimeEnd.IsZero() {
		sqler.AppendSql(" and tmr.update_time <= ?")
		sqler.AppendArg(updateTimeEnd)
	}
	sqler.SetAlias("tmr")
	return sqler
}
