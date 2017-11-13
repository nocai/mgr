package models

//
//import (
//	"mgr/util"
//	"github.com/astaxie/beego"
//)
//
//type Menu util.TreeNode
//
//type menuGenerator struct {
//	Ress     []Res
//	menuRess []Res
//	topRess  []ResVo
//}
//
//
//func  NewMenuGenerator(ress []Res) (*menuGenerator) {
//	return &menuGenerator{Ress:ress}
//}
//
//func (this *menuGenerator) Generate() []util.TreeNode {
//	this.findMenuRess()
//	this.assembleTrees()
//	return findMenuTrees(this.topRess)
//}
//
//func findMenuTrees(ress []ResVo) []util.TreeNode {
//	if len(ress) == 0 {
//		return make([]util.TreeNode, 0)
//	}
//
//	var trees []util.TreeNode
//	for _, res := range ress {
//		tree := util.TreeNode{Id:res.Id, Text:res.ResName, State:"open", Checked:false, Attributes:res.Path, Children:findMenuTrees(res.Children)}
//		trees = append(trees, tree)
//	}
//	return trees
//}
//
//func (this *menuGenerator) findMenuRess() {
//	if this.Ress == nil {
//		panic("ress is nil")
//	}
//
//	for _, res := range this.Ress {
//		if res.ResType == ResType_Menu {
//			this.menuRess = append(this.menuRess, res)
//		}
//	}
//}
//
//// 组装资源，形成树形结构
//func (this *menuGenerator) assembleTrees() {
//	for index, res := range this.menuRess {
//		if index == 0 {
//			topRes, err := getResVoTopByResId(res.Id)
//			if err != nil {
//				beego.Error(err)
//				return
//			}
//			this.topRess = append(this.topRess, *topRes)
//		} else {
//			this.assembleTree(&res)
//		}
//	}
//}
//
//func (this *menuGenerator) assembleTree(res *Res) {
//	for _, topRes := range this.topRess {
//		if topRes.Id == res.Id {
//			return
//		}
//		if assemble(&topRes, res) {
//			return
//		}
//	}
//	resVo := ResVo{Res:*res}
//	this.topRess = append(this.topRess, resVo)
//}
//
//func assemble(resVo *ResVo, res *Res) bool {
//	if res.Pid == resVo.Id {
//		resVo.Children = append(resVo.Children, ResVo{Res:*res})
//		return true
//	}
//	for _, child := range resVo.Children {
//		return assemble(&child, res)
//	}
//	return false;
//}
//
//func getResVoTopByResId(id int64) (*ResVo, error) {
//	res, err := GetResByResId(id)
//	if err != nil {
//		beego.Error(err)
//		return nil, ErrQuery
//	}
//	if res == nil {
//		return nil, nil
//	}
//
//	if isTop(res) {
//		return &ResVo{Res:*res}, nil
//	} else {
//		resVo, err := getResVoTopByResId(res.Pid)
//		if err != nil {
//			return nil, ErrQuery
//		}
//		resVo.Children = append(resVo.Children, ResVo{Res:*res})
//		return resVo, nil
//	}
//}
//
//func isTop(res *Res) bool {
//	if res.Pid == Pid_Default {
//		return true
//	}
//	return false
//}
//
