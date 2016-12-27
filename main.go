package main

import (
	"github.com/astaxie/beego"
	_ "mgr/routers"
	_ "mgr/models"
)

func main() {
	beego.SetStaticPath("/static", "static")
	beego.Run()
}
