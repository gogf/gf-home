package document

import (
	"github.com/gogf/gf-home/app/service/document"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

// 文档首页
func Index(r *ghttp.Request) {
	if r.IsAjaxRequest() {
		serveMarkdownAjax(r)
		return
	}
	path := r.GetString("path")
	if path == "" {
		if r.URL.RawQuery != "" {
			r.Response.RedirectTo("/index?" + r.URL.RawQuery)
		} else {
			r.Response.RedirectTo("/index")
		}
		return
	}
	config := g.Config()
	// 如果是静态文件请求，那么表示Web Server没有找到该文件，那么直接404，本接口不支持待后缀的静态文件处理。
	// 由于路由规则比较宽，这里也会有未存在的静态文件请求匹配进来。
	if r.IsFileRequest() {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	// 菜单内容
	var (
		title     = document.GetTitleByPath(path)
		baseTitle = config.GetString("document.title")
	)
	if title == "" {
		title = "404 NOT FOUND"
	}
	title += " - " + config.GetString("document.title")
	// markdown内容
	mdMainContent := document.GetMarkdown(path)
	mdMainContentParsed := document.ParseMarkdown(mdMainContent)
	r.Response.WriteTpl("document/index.html", g.Map{
		"title":               title,
		"baseTitle":           baseTitle,
		"mdMenuContentParsed": document.GetParsed("menus"),
		"mdMainContentParsed": mdMainContentParsed,
		"mdMainContent":       mdMainContent,
	})
}

// 文档更新hook
func UpdateHook(r *ghttp.Request) {
	if r.Get("password") == g.Config().GetString("document.hook") {
		document.UpdateDocGit()
		r.Response.Write("ok")
	} else {
		r.Response.WriteStatus(443)
	}
}

// 搜索文档
func Search(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"code": 1,
		"msg":  "",
		"data": document.SearchMdByKey(r.GetString("key")),
	})
}

// 处理ajax请求
func serveMarkdownAjax(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"code": 1,
		"msg":  "",
		"data": document.GetMarkdown(r.GetString("path", "index")),
	})
}
