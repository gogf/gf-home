package cdnreplace

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
)

// Middleware replaces the resource links with CDN links.
func Middleware(r *ghttp.Request) {
	r.Middleware.Next()
	cdnUrl := g.Cfg().GetString("cdn.url")
	if cdnUrl == "" {
		return
	}
	cdnVersion := g.Cfg().GetString("cdn.version")
	if cdnVersion == "" {
		cdnVersion = gdebug.BinVersion()
	}
	content := r.Response.BufferString()
	if len(content) > 0 {
		// HTML
		content, _ = gregex.ReplaceStringFuncMatch(`(href|src)=['"](.+?)['"]`, content, func(match []string) string {
			link := match[2]
			if len(link) == 0 {
				return match[0]
			}
			if link[0:1] != "/" && link[0:1] != "#" {
				if len(link) > 10 && link[0:10] == "javascript" {
					return match[0]
				}
				if len(link) > 7 && link[0:7] == "mailto:" {
					return match[0]
				}
				if len(link) > 4 && link[0:4] == "http" {
					return match[0]
				}
				link = "/" + link
			}
			if link[0:1] == "/" {
				switch gfile.ExtName(link) {
				case "png", "jpg", "jpeg", "gif", "svg", "bmp", "js", "css", "otf", "eot", "ttf", "woff", "woff2":
					return fmt.Sprintf(`%s="%s%s?%s"`, match[1], cdnUrl, link, cdnVersion)
				}
			}
			return match[0]
		})
		// Markdown
		content, _ = gregex.ReplaceStringFuncMatch(`!\[(.*?)\]\((.+?)\)`, content, func(match []string) string {
			link := match[2]
			if len(link) == 0 {
				return match[0]
			}
			if link[0:1] != "/" && link[0:1] != "#" {
				if len(link) > 10 && link[0:10] == "javascript" {
					return match[0]
				}
				if len(link) > 7 && link[0:7] == "mailto:" {
					return match[0]
				}
				if len(link) > 4 && link[0:4] == "http" {
					return match[0]
				}
				link = "/" + link
			}
			if link[0:1] == "/" {
				switch gfile.ExtName(link) {
				case "png", "jpg", "jpeg", "gif", "svg", "bmp":
					return fmt.Sprintf(`![%s](%s%s?%s)`, match[1], cdnUrl, link, cdnVersion)
				}
			}
			return match[0]
		})
		// Reset the response buffer.
		r.Response.SetBuffer([]byte(content))
	}
}
