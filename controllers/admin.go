package controllers

import (
	"mgr/conf"
	"mgr/models"
	"mgr/models/service/adminser"
	"mgr/util/key"

	"github.com/astaxie/beego"
	"mgr/models/service/userser"
)

type AdminController struct {
	BaseController
}

// Get method, for query
func (ctr *AdminController) Get() {
	page, _ := ctr.GetInt64("page", conf.Page)
	rows, _ := ctr.GetInt64("rows", conf.Rows)
	sort := ctr.GetString("sort")
	order := ctr.GetString("order")
	adminName := ctr.GetString("admin_name")
	//invalid, _ := ctr.GetInt("invalid", 0)
	if sort == "" {
		sort = "id"
		order = "asc"
	}

	key := key.New(page, rows, []string{sort}, []string{order})
	admin := &models.Admin{AdminName: "%" + adminName + "%"}
	pager, err := adminser.PageAdminVo(&adminser.AdminVoKey{Key: key, Admin: admin, Invalid: models.ValidAll})
	ctr.Print(pager, err)
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

	err := adminser.DeleteAdminById(id)
	ctr.PrintError(err)
}

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
