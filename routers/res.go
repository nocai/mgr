package routers

import (
	"github.com/astaxie/beego"
	"mgr/controllers"
)

func init() {
	beego.Router("/res", &controllers.ResController{}) // Get method
	beego.Router("/res/:id:int", &controllers.ResController{}) // Get method

	beego.Router("/resSelects", &controllers.ResSelectController{}) // Get method
}