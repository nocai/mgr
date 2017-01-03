package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"errors"
	"time"
	"sync"
)

var (
	ErrResNameExist = errors.New("资源名称存在")
)

func insertRes(ch chan error, o orm.Ormer, res *Res) {
	//defer mutex.Unlock()
	var err error = nil

	mutex := &sync.Mutex{}
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		ch <- err
	}()

	if res.ResName == "" {
		err = ErrArgument
		return
	}

	if existOfResName(res.Id, res.ResName) {
		err = ErrResNameExist
		return
	}

	if res.Path == "" {
		err = ErrArgument
		return
	}
	fillResTime(res)

	id, err := o.Insert(res);
	if err != nil {
		beego.Error(err)
		err = ErrInsert
		return
	}

	beego.Debug(fmt.Sprintf("id = %v", id))
	return
}

func InsertRes(res *Res) (error) {
	o := orm.NewOrm()
	o.Begin()

	chs := make([]chan error, len(res.Children) + 1)
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan error)
		if i == 0 {
			go insertRes(chs[i], o, res)
		} else {
			go insertRes(chs[i], o, &res.Children[i])
		}
	}

	for _, ch := range chs {
		if <-ch != nil {
			o.Rollback()
			return ErrInsert
		}
	}

	o.Commit()
	return nil
}

func fillResTime(res *Res) {
	now := time.Now()
	res.CreateTime = now
	res.UpdateTime = now
}

func existOfResName(id int64, resName string) bool {
	res, err := GetResByResName(resName)
	if err == nil {
		if res.Id != id {
			beego.Debug(fmt.Sprintf("资源名称[%v]存在", resName))
			return true
		}
	}
	return false
}

func GetResByResName(resName string) (*Res, error) {
	res := &Res{ResName:resName}

	if err := orm.NewOrm().Read(res, "ResName"); err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}
	return res, nil
}