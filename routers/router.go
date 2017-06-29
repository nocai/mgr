package routers

import (
	"github.com/astaxie/beego"
	"mgr/controllers"
	"mgr/routers/filter"
)

func init() {
	beego.InsertFilter("/*",beego.BeforeExec,filter.FilterUser)
	//
	beego.Router("/t", &controllers.MainController{}, "get:T")
	beego.Router("/", &controllers.MainController{})
	//
	//beego.Router("/login", &controllers.LoginController{})
	//
	//beego.Router("/nav", &controllers.NavController{})

	beego.Router("/:path.html", &controllers.HtmlController{})

}

// admins
func init() {
	beego.Router("/admins", &controllers.AdminController{})
	beego.Router("/admins/:id:int", &controllers.AdminController{})

	beego.Router("/adminValid", &controllers.AdminInvalidController{})
	beego.Router("/adminValid/:id:int/:invalid:int", &controllers.AdminInvalidController{})
}

// roles
func init() {
	beego.Router("/roles", &controllers.RoleController{})
	beego.Router("/roles/:id:int", &controllers.RoleController{})
}

// res
func init() {
	//beego.Router("/res", &controllers.ResController{}) // Get method
	//beego.Router("/res/:id:int", &controllers.ResController{}) // Get method
	//beego.Router("/res/:id:int", &controllers.ResController{}) // Delete method
	//
	//beego.Router("/resSelects", &controllers.ResSelectController{}) // Get method
}