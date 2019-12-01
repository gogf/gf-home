package router

import (
	"github.com/gogf/gf-home/app/api/cli/binary"
	"github.com/gogf/gf-home/app/api/cli/project"
	"github.com/gogf/gf-home/app/api/cli/version"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	g.Server().Group("/cli", func(g *ghttp.RouterGroup) {
		g.ALL("/*path", binary.Index)
		g.ALL("/version", version.Latest)
		g.Group("/binary", func(g *ghttp.RouterGroup) {
			g.ALL("/*path", binary.Index)
			g.ALL("/md5", binary.Md5)
		})
		g.Group("/project", func(g *ghttp.RouterGroup) {
			g.ALL("/md5", project.Md5)
			g.ALL("/zip", project.Zip)
		})
	})
}
