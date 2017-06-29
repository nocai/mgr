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
