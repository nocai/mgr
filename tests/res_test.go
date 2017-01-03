package test

import (
	"testing"
	"mgr/models"
	"mgr/util"
)

func TestRes(t *testing.T) {
	resName := string(util.Krand(6, util.KC_RAND_KIND_UPPER))
	res := &models.Res{ResName:resName}
	models.InsertRes(res)
}
