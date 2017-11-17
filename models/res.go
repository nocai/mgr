package models

//
//import (
//	"github.com/astaxie/beego/orm"
//	"github.com/astaxie/beego"
//	"fmt"
//	"errors"
//	"time"
//	"sync"
//	"mgr/util"
//	"sort"
//)
//

//
//const Pid_Default = -1
//
//var (
//	ErrResNameExist = errors.New("资源名称存在")
//)
//

//
//
//type ResSorter []ResVo
//
//func NewResSorter(ress []ResVo) ResSorter {
//	rs := make(ResSorter, 0, len(ress))
//
//	for _, v := range ress {
//		rs = append(rs, v)
//	}
//
//	return rs
//}
//
//func (rs ResSorter) getResVos() []ResVo {
//	ress := []ResVo{}
//	for _, res := range rs {
//		ress = append(ress, res)
//	}
//	return ress
//}
//
//func (rs ResSorter) Len() int {
//	return len(rs)
//}
//
//func (rs ResSorter) Less(i, j int) bool {
//	return rs[i].Seq < rs[j].Seq
//}
//
//func (rs ResSorter) Swap(i, j int) {
//	rs[i], rs[j] = rs[j], rs[i]
//}
//

//
//type ResVo struct {
//	Res
//	ResTypeDesc string
//
//	Children    []ResVo
//}
//
//func FindResVoByKey(key *ResKey, cascade bool) ([]ResVo, error) {
//	ress, err := FindResByKey(key)
//
//	if err != nil || len(ress) == 0 {
//		return make([]ResVo, 0), err
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(len(ress))
//
//	var resVos []ResVo
//	for _, res := range ress {
//		go func(res Res) {
//			defer wg.Done()
//
//			resVo := ResVo{Res:res}
//			if cascade {
//				children, err := FindResVoByKey(&ResKey{Res:Res{Pid:res.Id}}, false)
//				if err != nil {
//					resVo.Children = make([]ResVo, 0)
//				} else {
//					resVo.Children = children
//				}
//			} else {
//				resVo.Children = make([]ResVo, 0)
//			}
//			resVos = append(resVos, resVo)
//		}(res)
//	}
//	wg.Wait()
//
//	rs := NewResSorter(resVos)
//	sort.Sort(rs)
//	return rs.getResVos(), nil
//}