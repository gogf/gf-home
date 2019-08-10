package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	// 管理接口
	g.Server().EnableAdmin("/admin")

	// 某些浏览器会直接请求/favicon.ico文件，会产生404
	g.Server().SetRewrite("/favicon.ico", "/resource/image/favicon.ico")

	// 为平滑重启管理页面设置HTTP Basic账号密码
	g.Server().BindHookHandler("/admin/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
		user := g.Config().GetString("admin.user")
		pass := g.Config().GetString("admin.pass")
		if !r.BasicAuth(user, pass) {
			r.ExitAll()
		}
	})
}
