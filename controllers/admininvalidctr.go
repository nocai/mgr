package controllers

import (
	"mgr/models"
	"mgr/models/service/adminser"
	"mgr/models/service/userser"
)

type AdminInvalidController struct {
	BaseController
}

func (aic *AdminInvalidController) Get() {
	m := adminser.FindAdminValids()
	aic.PrintData(m)
}

func (aic *AdminInvalidController) Put() {
	id, _ := aic.GetInt(":id")
	invalid, _ := aic.GetInt(":invalid", 0)

	admin, err := adminser.GetAdminById(int64(id))
	if err != nil {
		aic.PrintError(err)
		return
	}
	user, err := userser.GetUserById(admin.UserId)
	if err != nil {
		aic.PrintError(err)
		return
	}
	user.Invalid = models.ValidEnum(invalid)
	err = userser.UpdateUser(user)
	aic.PrintError(err)
}
