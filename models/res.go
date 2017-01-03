package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"errors"
	"time"
	"sync"
	"mgr/util"
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

// res.Children的数据同样会被入库,对应的Id会被填充
func InsertRes(res *Res) (error) {
	o := orm.NewOrm()
	o.Begin()

	ch := make(chan error)
	go insertRes(ch, o, res)
	if <- ch != nil {
		o.Rollback()
		return ErrInsert
	}

	if len(res.Children) > 0 {
		chs := make([]chan error, len(res.Children))
		for i := 0; i < len(res.Children); i++ {
			chs[i] = make(chan error)
			child := res.Children[i]
			if child.Pid == 0 {
				child.Pid = res.Id
			}
			go insertRes(chs[i], o, &child)
		}

		for _, ch := range chs {
			if <-ch != nil {
				o.Rollback()
				return ErrInsert
			}
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

func PageRes(key *util.PagerKey) (*util.Pager, error) {
	key.AppendDataSql("select * from t_mgr_res as tmr where 1 = 1")

	if resName, ok := key.Data["resName"].(string); ok && resName != "" {
		key.AppendDataSql(" and tmr.res_name like ?")
		key.AppendArg("%" + resName + "%")
	}

	return nil, nil

}