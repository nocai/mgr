package controllers

import (
	"mgr/conf"
	"mgr/models"
	"mgr/models/service/adminser"
	"mgr/util/key"

	"github.com/astaxie/beego"
)

type AdminController struct {
	BaseController
}


// Get method, for query
func (ctr *AdminController) Get() {
	defer ctr.recoverPanic()
	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort", "id")
	order := ctr.GetString("order", "asc")
	adminName := ctr.GetString("admin_name")
	invalid, _ := ctr.GetInt("invalid", 2)

	key := key.New(page, rows, []string{sort}, []string{order})
	admin := &models.Admin{AdminName: "%" + adminName + "%"}
	pager := adminser.PageAdminVo(&adminser.AdminVoKey{Key: key, Admin: admin, Invalid: models.ValidEnum(invalid)})
	ctr.PrintData(pager)
}

// 添加 修改
func (ctr *AdminController) Post() {
	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = ", id)
	adminName := ctr.GetString("admin_name")
	password := ctr.GetString("password")

	if id == 0 {
		// 添加
		adminVo := &adminser.AdminVo{
			Admin: &models.Admin{
				AdminName: adminName,
			},
			User: &models.User{
				Username: adminName,
				Password: password,
				Invalid:  models.Invalid,
			},
		}
		ctr.PrintError(adminser.InsertAdminVo(adminVo))
	} else {
		// 更新
		adminVo, err := adminser.GetAdminVoById(id)
		if err != nil {
			beego.Error(err)
			ctr.PrintError(err)
			return
		}
		adminVo.AdminName = adminName
		adminVo.User.Username = adminName
		adminVo.User.Password = password
		ctr.PrintError(adminser.UpdateAdminVo(adminVo))
	}
}

func (ctr *AdminController) Delete() {
	beego.Debug(ctr.Input())

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = ", id)

	adminser.DeleteAdminById(id)
	ctr.PrintOk("删除成功")
}

