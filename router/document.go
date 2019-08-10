package router

import (
	"github.com/gogf/gf-home/app/api/document"
	"github.com/gogf/gf/frame/g"
)

func init() {
	g.Server().BindHandler("/*path", document.Index)
	g.Server().BindHandler("/hook", document.UpdateHook)
	g.Server().BindHandler("/search", document.Search)
}
