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
	ctr.debugInput()

	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")
	adminName := ctr.GetString("admin_name")
	valid, _ := ctr.GetInt("valid", 0)
	if sort == "" {
		sort = "id"
		order = "asc"
	}

	key := key.New(page, rows, []string{sort}, []string{order}, true)
	admin := &models.Admin{AdminName: "%" + adminName + "%", Invalid:models.ValidEnum(valid)}
	pager, err := adminser.PageAdmin(&models.AdminKey{Key: key, Admin: admin})
	ctr.PrintResult(pager, err)
}

// 添加 修改
func (ctr *AdminController) Post() {
	ctr.debugInput()

	id, _ := ctr.GetInt64(":id", 0)
	beego.Debug("id = ", id)
	adminName := ctr.GetString("admin_name")
	password := ctr.GetString("password")

	if id == 0 {
		// 添加
		adminVo := &models.AdminVo{
			Admin:  &models.Admin{
				AdminName: adminName,
			},
			User: &models.User{
				Username: adminName,
				Password: password,
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
	beego.Debug("id = %v", id)

	err := adminser.DeleteAdminById(id)
	ctr.PrintError(err)
}
