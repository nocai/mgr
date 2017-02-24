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

type Res struct {
	Id         int64 `json:"id"`
	ResName    string `json:"res_name"`
	Path       string `json:"path"`
	Level      int `json:"level"`

	Pid        int64 `json:"pid"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`

	Children   []Res `orm:"-"`
}

// 多字段唯一键
func (res *Res) TableUnique() [][]string {
	return [][]string{
		[]string{"ResName"},
	}
}

// 多字段索引
func (res *Res) TableIndex() [][]string {
	return [][]string{
		[]string{"ResName"},
	}
}

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
		beego.Error("参数无效")
		err = ErrArgument
		return
	}

	if existOfResName(res.Id, res.ResName) {
		err = ErrResNameExist
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
func InsertRes(resVo *ResVo) (error) {
	o := orm.NewOrm()
	o.Begin()

	ch := make(chan error)
	go insertRes(ch, o, &resVo.Res)
	if <-ch != nil {
		o.Rollback()
		return ErrInsert
	}

	if len(resVo.Children) > 0 {
		chs := make([]chan error, len(resVo.Children))
		for i := 0; i < len(resVo.Children); i++ {
			chs[i] = make(chan error)
			if resVo.Children[i].Pid == 0 {
				resVo.Children[i].Pid = resVo.Id
			}
			go insertRes(chs[i], o, &resVo.Children[i].Res)
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

func UpdateRes(res *Res) error {
	if existOfResName(res.Id, res.ResName) {
		return ErrResNameExist
	}
	res.UpdateTime = time.Now()

	o := orm.NewOrm()

	affected, err := o.Update(res)
	if err != nil {
		beego.Error(err)
		return ErrUpdate
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
	res, err := GetResByResName(resName)
	beego.Info(fmt.Sprintf("%+v", res))
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

func PageRes(key *ResKey) (*util.Pager, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var total int64
	err := o.Raw(sqler.GetCountSql(), sqler.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, 0, make([]Res, 0)), ErrQuery
	}
	if total == 0 {
		return util.NewPager(key.Key, 0, make([]Res, 0)), nil
	}

	var res []Res
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&res)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key.Key, total, make([]Res, 0)), ErrQuery
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return util.NewPager(key.Key, total, res), nil
}

type ResKey struct {
	*util.Key

	Res

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time

	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
}

func (this *ResKey) getSqler() *util.Sqler {
	sqler := &util.Sqler{Key:this.Key}

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
	if level := this.Level; level != 0 {
		sqler.AppendSql(" and tmr.level = ?")
		sqler.AppendArg(level)
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
	return sqler

}

func FindResByKey(key *ResKey) ([]Res, error) {
	o := orm.NewOrm()
	sqler := key.getSqler()

	var ress []Res
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&ress)
	if err != nil {
		beego.Debug(err)
		return make([]Res, 0), ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return ress, nil
}

func DeleteResById(id int64) error {
	if id == 0 {
		return nil
	}

	key := &ResKey{Res:Res{Pid:id}}
	ress, err := FindResByKey(key)
	if err == nil && len(ress) > 0 {
		return errors.New("请先删除子资源")
	}

	affected, err := orm.NewOrm().Delete(&Res{Id:id})
	if err != nil {
		beego.Error(err)
		return ErrDelete
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}

type ResVo struct {
	Res
	Children []ResVo
}

func FindResVoByKey(key *ResKey, cascade bool) ([]ResVo, error) {
	ress, err := FindResByKey(key)

	if err != nil || len(ress) == 0 {
		return make([]ResVo, 0), err
	}

	var wg sync.WaitGroup
	wg.Add(len(ress))

	var resVos []ResVo
	for _, res := range ress {
		go func(res Res) {
			defer wg.Done()

			resVo := ResVo{Res:res}
			if cascade {
				children, err := FindResVoByKey(&ResKey{Res:Res{Pid:res.Id}}, false)
				if err != nil {
					resVo.Children = make([]ResVo, 0)
				} else {
					resVo.Children = children
				}
			} else {
				resVo.Children = make([]ResVo, 0)
			}
			resVos = append(resVos, resVo)
		}(res)
	}
	wg.Wait()

	return resVos, nil
}