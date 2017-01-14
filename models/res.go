package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"errors"
	"time"
	"sync"
	"mgr/util"
	"bytes"
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

type ResKey struct {
	Res
}


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

	//Children   []Res `orm:"-"`
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

func PageRes(key *util.PagerKey) (*util.Pager, error) {
	key.AppendSql("select * from t_mgr_res as tmr where 1 = 1")

	if resName, ok := key.Data["resName"].(string); ok && resName != "" {
		key.AppendSql(" and tmr.res_name like ?")
		key.AppendArg("%" + resName + "%")
	}
	if path, ok := key.Data["path"].(string); ok && path != "" {
		key.AppendSql(" and tmr.path like = ?")
		key.AppendArg("%" + path + "%")
	}
	if pid, ok := key.Data["pid"].(int64); ok {
		key.AppendSql(" and tmr.pid = ?")
		key.AppendArg(pid)
	}

	o := orm.NewOrm()

	var total int64
	err := o.Raw(key.GetCountSql(), key.GetArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key, 0, make([]Res, 0)), ErrQuery
	}
	if total == 0 {
		return util.NewPager(key, 0, make([]Res, 0)), nil
	}

	var res []Res
	affected, err := o.Raw(key.GetSql(), key.GetArgs()).QueryRows(&res)
	if err != nil {
		beego.Error(err)
		return util.NewPager(key, total, make([]Res, 0)), ErrQuery
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return util.NewPager(key, total, res), nil
}

type ResKey struct {
	util.Key
	Res

	CreateTimeStart time.Time
	CreateTimeEnd   time.Time

	UpdateTimeStart time.Time
	UpdateTimeEnd   time.Time
}

func FindResByKey(key *ResKey) ([]Res, error) {
	key.AppendSql(`select * from t_mgr_res as tmr where 1 = 1`)

	if key.Id != 0 {
		key.AppendSql(" and tmr.id = ?")
		key.AppendArg(key.Id)
	}
	if key.ResName != "" {
		key.AppendSql(" and tmr.res_name = ?")
		key.AppendArg(key.ResName)
	}
	if key.Path != "" {
		key.AppendSql(" and tmr.path = ?")
		key.AppendArg(key.Path)
	}
	if level := key.Level; level != 0 {
		key.AppendSql(" and tmr.level = ?")
		key.AppendArg(level)
	}
	if pid := key.Pid; pid != 0 {
		key.AppendSql(" and tmr.pid = ?")
		key.AppendArg(pid)
	}
	if createTime := key.CreateTime; !createTime.IsZero() {
		key.AppendSql(" and tmr.create_time = ?")
		key.AppendArg(createTime)
	}
	if updateTime := key.UpdateTime; !updateTime.IsZero() {
		key.AppendSql(" and tmr.update_time = ?")
		key.AppendArg(updateTime)
	}
	if createTimeStart := key.CreateTimeStart; !createTimeStart.IsZero() {
		key.AppendSql(" and tmr.create_time >= ?")
		key.AppendArg(createTimeStart)
	}
	if createTimeEnd := key.CreateTimeEnd; !createTimeEnd.IsZero() {
		key.AppendSql(" and tmr.create_time <= ?")
		key.AppendArg(createTimeEnd)
	}
	if updateTimeStart := key.UpdateTimeStart; !updateTimeStart.IsZero() {
		key.AppendSql(" and tmr.update_time >= ?")
		key.AppendArg(updateTimeStart)
	}
	if updateTimeEnd := key.UpdateTimeStart; !updateTimeEnd.IsZero() {
		key.AppendSql(" and tmr.update_time <= ?")
		key.AppendArg(updateTimeEnd)
	}

	o := orm.NewOrm()

	var ress []Res
	affected, err := o.Raw(key.GetSql(), key.GetArgs()).QueryRows(&ress)
	if err != nil {
		beego.Debug(err)
		return make([]Res, 0), ErrQuery
	}
	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return ress, nil
}

func FindResByKey(key *ResKey) ([]Res , error) {
	o := orm.NewOrm()

	var buf bytes.Buffer
	buf.WriteString(`select * from t_mgr_res as tmr where 1 = 1`)
	if key.Id != 0 {
		buf.WriteString(` and tmr.id = ?`)
	}

	sql := `select * from t_mgr_res where 1 = 1`
	bytes.Buffer{}

	var ress []Res

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
				resVo.Children =  make([]ResVo,0)
			}
			resVos = append(resVos, resVo)
		}(res)
	}
	wg.Wait()

	return resVos, nil
}