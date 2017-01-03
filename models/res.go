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
		close(ch)
		return ErrInsert
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	ch <- nil
	return nil
}

func InsertRes(res *Res) (error) {
	children := res.children
	children = append(children, *res)
	fmt.Println(len(children))
	o := orm.NewOrm()
	o.Begin()
	ch := make(chan error, len(children))
	for _, child := range children {
		go insertRes(o, ch, &child)
	}

	close(ch)
	for err := range ch {
		if err != nil {
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