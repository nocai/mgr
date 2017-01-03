package test

import (
	"testing"
	"mgr/models"
	"mgr/util"
	"github.com/astaxie/beego"
	"fmt"
)

func TestRes(t *testing.T) {
	resName := string(util.Krand(6, util.KC_RAND_KIND_UPPER))
	res := &models.Res{ResName:resName, Path:"aaaaa"}

	//res.Children = append(res.Children, models.Res{ResName:string(util.Krand(8, util.KC_RAND_KIND_UPPER)), Path:"path"})

	//beego.Error(res)
	err := models.InsertRes(res)
	if err != nil {
		beego.Error(err)
		t.FailNow()
	}
	beego.Info(fmt.Sprintf("%+v", res))
}
