package routers

import (
	"mgr/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/admins", &controllers.AdminController{})
	beego.Router("/admins/:id:int", &controllers.AdminController{})
}
