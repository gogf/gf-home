package adminauth

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// Middleware does the authentication for administration.
func Middleware(r *ghttp.Request) {
	user := g.Config().GetString("admin.user")
	pass := g.Config().GetString("admin.pass")
	if !r.BasicAuth(user, pass) {
		return
	}
	r.Middleware.Next()
}
