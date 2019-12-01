package router

import (
	"github.com/gogf/gf-home/app/service/middleware/adminauth"
	"github.com/gogf/gf/frame/g"
)

func init() {
	// Administration interface.
	g.Server().EnableAdmin("/admin")

	// Basic authentication for admin.
	g.Server().BindMiddleware("/admin/*", adminauth.Middleware)

	// Avoid 404 for some browsers requesting "/favicon.ico".
	g.Server().SetRewrite("/favicon.ico", "/resource/image/favicon.ico")
}
