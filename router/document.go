package router

import (
	"github.com/gogf/gf-home/app/api/document"
	"github.com/gogf/gf-home/app/service/middleware/cdnreplace"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	// Documentation service.
	g.Server().Group("/", func(g *ghttp.RouterGroup) {
		g.Middleware(cdnreplace.Middleware)
		g.ALL("/*path", document.Index)
		g.ALL("/hook", document.UpdateHook)
		g.ALL("/search", document.Search)
	})
}
