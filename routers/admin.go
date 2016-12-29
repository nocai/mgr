package routers

import (
	"github.com/astaxie/beego"
	"mgr/controllers"
)

func init() {
	beego.Router("/admins", &controllers.AdminController{})
	beego.Router("/admins/:id:int", &controllers.AdminController{})
}
