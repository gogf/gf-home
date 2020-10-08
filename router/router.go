package router

import (
	"github.com/gogf/gf/frame/g"
)

func init() {
	// Avoid 404 for some browsers requesting "/favicon.ico".
	g.Server().SetRewrite("/favicon.ico", "/resource/image/favicon.ico")
}
