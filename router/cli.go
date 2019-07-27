package router

import (
	"github.com/gogf/gf-home/app/api/cli/binary"
	"github.com/gogf/gf-home/app/api/cli/project"
	"github.com/gogf/gf-home/app/api/cli/version"
	"github.com/gogf/gf/g"
)

func init() {
	g.Server().BindHandler("/cli/*path", binary.Index)
	g.Server().BindHandler("/cli/binary/*path", binary.Index)
	g.Server().BindHandler("/cli/binary/md5", binary.Md5)

	g.Server().BindHandler("/cli/project/md5", project.Md5)
	g.Server().BindHandler("/cli/project/zip", project.Zip)

	g.Server().BindHandler("/cli/version", version.Latest)
}
