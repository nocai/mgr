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

func insertRes(mutex *sync.Mutex, ch chan error, o orm.Ormer, res *Res) {
	//defer mutex.Unlock()
	var err error = nil
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
	fillTime(res)

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

	var resList []Res
	resList = append(resList, *res)
	if len(res.Children) > 0 {
		resList = append(resList, res.Children...)
	}
	beego.Error(fmt.Sprintf("%+v", resList))
	chs := make([]chan error, len(resList))
	mutex := &sync.Mutex{}
	for i, r := range resList {
		chs[i] = make(chan error)
		go insertRes(mutex, chs[i], o, &r)
	}

	for _, ch := range chs {
		err := <-ch
		fmt.Println(err)
		if err != nil {
			beego.Error(err)
			o.Rollback()
			return ErrInsert
		}
	}
	o.Commit()
	return nil
}

func fillTime(res *Res) {
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

	if err := orm.NewOrm().Read(res); err != nil {
		beego.Error(err)
		return nil, ErrQuery
	}
	return res, nil
}