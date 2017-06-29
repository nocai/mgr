package routers

import (
	"github.com/astaxie/beego"
	"mgr/controllers"
)

func init() {

	// roles
	beego.Router("/roles", &controllers.RoleController{})
	beego.Router("/roles/:id:int", &controllers.RoleController{})
}
