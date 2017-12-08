package controllers

import (
	"mgr/models/service/adminser"
	"mgr/models/service/userser"
)

type AdminInvalidController struct {
	BaseController
}

func (aic *AdminInvalidController) Get() {
	m := adminser.FindAdminValids()
	aic.PrintJson(m)
}

func (aic *AdminInvalidController) Put() {
	id, _ := aic.GetInt(":id")
	invalid, _ := aic.GetInt(":invalid", 0)

	admin, err := adminser.GetAdminById(int64(id))
	if err != nil {
		aic.PrintJson(err)
		return
	}
	user, err := userser.GetUserById(admin.UserId)
	if err != nil {
		aic.PrintJson(err)
		return
	}
	user.Invalid = userser.ValidEnum(invalid)
	err = userser.UpdateUser(user)
	aic.PrintJson(err)
}
