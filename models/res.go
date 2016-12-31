package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"errors"
	"time"
)

var (
	ErrResNameExist = errors.New("资源名称存在")
)

func insertRes(o orm.Ormer, ch chan error, res *Res) error {
	if res.ResName == "" {
		return ErrArgument
	}
	if existOfResName(res.Id, res.ResName) {
		return ErrResNameExist
	}

	if res.Path == "" {
		return ErrArgument
	}
	fillTime(res)

	affected, err := o.Insert(res);
	if err != nil {
		beego.Error(err)
		o.Rollback()
		ch <- err
		return ErrInsert
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	ch <- nil
	return nil
}

func InsertRes(res *Res) (error) {
	o := orm.NewOrm()
	o.Begin()

	var ch chan error
	if len(res.children) == 0 {
		go insertRes(o, ch, res)
	} else {
		res.children = append(&res)
		ch := make(chan error, len(res.children))
		for _, res := range res.children {
			go insertRes(o, ch, &res)
		}
	}

	if len(res.children) > 0 {

		for err := range ch {
			if err != nil {
				return ErrInsert
			}
		}
	}
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