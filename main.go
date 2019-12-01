package main

import (
	_ "github.com/gogf/gf-home/boot"
	_ "github.com/gogf/gf-home/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
