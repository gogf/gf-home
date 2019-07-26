package router

import (
	"github.com/gogf/gf-home/app/api/project"
	"github.com/gogf/gf/g"
)

func init() {
	g.Server().BindHandler("/project/md5", project.Md5)
	g.Server().BindHandler("/project/zip", project.Zip)
}
