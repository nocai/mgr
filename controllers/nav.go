package controllers

import (
	"mgr/models"
	"github.com/astaxie/beego"
	"mgr/util"
)

type NavController struct {
	BaseController
}

func (ctr *NavController) Get() {
	key := &models.ResKey{Res:models.Res{Pid:-1}}
	ress, err := models.FindResVoByKey(key, true)
	if err != nil {
		beego.Error(err)
	}

	trees := toTrees(ress)
	ctr.Print(trees)
}

func toTrees(ress []models.ResVo) []util.TreeNode {
	if len(ress) == 0 {
		return make([]util.TreeNode, 0)
	}

	var trees []util.TreeNode
	for _, res := range ress {
		tree := util.TreeNode{Id:res.Id, Text:res.ResName, State:"open", Checked:false, Attributes:res.Path, Children:toTrees(res.Children)}
		trees = append(trees, tree)
	}
	return trees
}

