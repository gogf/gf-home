package router

import (
	"github.com/gogf/gf-home/app/api/document"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	g.Server().Group("/", func(g *ghttp.RouterGroup) {
		g.ALL("/*path", document.Index)
		g.ALL("/hook", document.UpdateHook)
		g.ALL("/search", document.Search)
	})
}
