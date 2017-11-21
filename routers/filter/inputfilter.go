package filter

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var FilterInput = func(ctx *context.Context) {
	fmt.Println("FilterInput...")

	if beego.BConfig.RunMode == beego.DEV {
		beego.Info("请求路径:", fmt.Sprintf("%#v", ctx.Input.Params()))
		beego.Info("输入参数:", ctx.Input)
	}
}
