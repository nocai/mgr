package filter

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

var FilterUser = func(ctx *context.Context) {
	fmt.Println("filter...")

}
