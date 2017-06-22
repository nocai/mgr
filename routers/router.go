package routers

import (
	"github.com/astaxie/beego"
	"mgr/controllers"
)

func init() {
	//beego.InsertFilter("/*",beego.BeforeExec,filter.FilterUser)
	//
	beego.Router("/t", &controllers.MainController{}, "get:T")
	beego.Router("/", &controllers.MainController{})
	//
	//beego.Router("/login", &controllers.LoginController{})
	//
	//beego.Router("/nav", &controllers.NavController{})

	// roles
	beego.Router("/roles", &controllers.RoleController{})
	beego.Router("/roles/:id:int", &controllers.RoleController{})
}
