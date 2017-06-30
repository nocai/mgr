package main

import (
	"github.com/astaxie/beego"
	_ "mgr/models"
	_ "mgr/routers"
	"net/http"
	"html/template"
	"fmt"
)

func main() {
	beego.SetStaticPath("/static", "static")
	beego.ErrorHandler("dbError",dbError)
	beego.Run()
}
func dbError(rw http.ResponseWriter, r *http.Request){
	fmt.Println(beego.BConfig.WebConfig.ViewsPath)
	t,_:= template.New("dberror.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath+"/dberror.html")
	data :=make(map[string]interface{})
	data["content"] = "database is now down"
	t.Execute(rw, data)
}