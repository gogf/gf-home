package boot

import (
	"github.com/gogf/gf/frame/g"
)

// 用于配置初始化.
func init() {
	c := g.Config()
	s := g.Server()

	s.AddSearchPath(c.GetString("document.path"))
}
