package test

import (
	"github.com/astaxie/beego"
	"mgr/models"
	"testing"
	"fmt"
)

func TestNav(t *testing.T) {
	key := &models.ResKey{Res:models.Res{Pid:-1}}
	ress, err := models.FindResVoByKey(key, true)
	if err != nil {
		beego.Error(err)
	}

	beego.Debug(fmt.Sprintf("%+v", ress))
}
