package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
	_ "mgr/models"
	_ "mgr/routers"
	"net/http"
)

func main() {
	beego.SetStaticPath("/static", "static")
	beego.ErrorHandler("dbError", dbError)
	beego.Run()
}

func dbError(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(beego.BConfig.WebConfig.ViewsPath)
	t, _ := template.New("dberror.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/dberror.html")
	data := make(map[string]interface{})
	data["content"] = "database is now down"
	t.Execute(rw, data)
}
