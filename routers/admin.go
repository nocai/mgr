package routers

import (
	"github.com/astaxie/beego"
	"mgr/controllers"
)

func init() {
	beego.Router("/admins", &controllers.AdminController{})
}
