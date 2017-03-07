package filter

import (
	"github.com/astaxie/beego/context"
	"fmt"
)

var FilterUser = func(ctx *context.Context) {
	fmt.Println("filter...")
}

