package router

import (
	"github.com/gogf/gf-home/app/api/cli"
	"github.com/gogf/gf/g"
)

func init() {
	g.Server().BindHandler("/cli/*path", cli.Index)
	g.Server().BindHandler("/cli/download", cli.Download)
}
