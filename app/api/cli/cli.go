package cli

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

var (
	cliRoot = g.Config().GetString("cli.path")
)

// CLI二进制文件浏览
func Index(r *ghttp.Request) {
	if r.Get("path") == "" {
		r.Response.WriteStatus(404)
		r.Exit()
	}
	r.Response.ServeFile(cliRoot+"/"+r.Get("path"), true)
}

// CLI二进制文件下载
func Download(r *ghttp.Request) {
	os := r.Get("os")
	arch := r.Get("arch")
	name := "gf"
	if os == "windows" {
		name += ".exe"
	}
	r.Response.ServeFileDownload(
		cliRoot+"/"+r.Get("path")+fmt.Sprintf(`latest/%s_%s/%s`, os, arch, name),
		name,
	)
}
