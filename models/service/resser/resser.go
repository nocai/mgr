package resser

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/models/service"
	"mgr/util/pager"
	"time"
	"mgr/conf"
)

// 分页
func PageRes(key *ResKey) *pager.Pager {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var total int64
	err := o.Raw(sqler.GetCountSql(), sqler.GetArgs()).QueryRow(&total)
	if err != nil {
		panic(service.NewError(conf.MsgQuery, err))
	}

	ress := FindResByKey(key)
	return pager.New(key.Key, total, ress)
}

// Query
func FindResByKey(key *ResKey) []Res {
	o := orm.NewOrm()
	sqler := key.getSqler()

	ress := make([]Res, 0)
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&ress)
	if err != nil {
		panic(service.NewError(conf.MsgQuery, err))
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return ress
}

func InsertRes(res *Res) error {
	if existOfResName(res.Id, res.ResName) {
		return ErrResNameExist
	}
	now := time.Now()
	if res.CreateTime.IsZero() {
		res.CreateTime = now
	}
	res.UpdateTime = now

	o := orm.NewOrm()
	resId, err := o.Insert(res)
	if err != nil {
		panic(service.NewError(conf.MsgInsert, err))
	}
	res.Id = resId
	return nil
}

func UpdateRes(res *Res) error {
	exist := existOfResName(res.Id, res.ResName)
	if exist {
		return ErrResNameExist
	}
	res.UpdateTime = time.Now()

	o := orm.NewOrm()

	affected, err := o.Update(res)
	if err != nil {
		panic(service.NewError(conf.MsgUpdate, err))
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}

func fillResTime(res *Res) {
	now := time.Now()
	res.CreateTime = now
	res.UpdateTime = now
}

func existOfResName(id int64, resName string) bool {
	res := GetResByResName(resName)
	if res != nil && res.Id != id {
		beego.Debug(fmt.Sprintf("资源名称[%v]存在", resName))
		return true
	}
	return false
}

func GetResByResId(id int64) *Res {
	key := &ResKey{Res: Res{Id: id}}
	ress := FindResByKey(key)
	if len(ress) == 0 {
		return nil
	}
	if len(ress) > 1 {
		beego.Warn(service.ErrDataDuplication, "id = ", id)
	}
	return &ress[0]
}

func GetResByResName(resName string) *Res {
	key := &ResKey{Res: Res{ResName: resName}}
	ress := FindResByKey(key)
	if len(ress) == 0 {
		return nil
	}
	if len(ress) > 1 {
		beego.Warn(service.ErrDataDuplication, "resName = ", resName)
	}
	return &ress[0]
}

func DeleteResById(id int64) error {
	if id <= 0 {
		return service.ErrArgument
	}

	key := &ResKey{Res: Res{Pid: id}}
	ress := FindResByKey(key)
	if len(ress) > 0 {
		return errors.New("请先删除子资源")
	}

	affected, err := orm.NewOrm().Delete(&Res{Id: id})
	if err != nil {
		beego.Error(err)
		return service.NewError(conf.MsgDelete, err)
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}
