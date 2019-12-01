package version

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/text/gregex"
	"time"
)

var (
	githubTagUrl  = g.Config().GetString("git.github")
	latestVersion = g.NewVar(nil, true)
)

func init() {
	if githubTagUrl == "" {
		glog.Fatal("Git configuration cannot be empty")
	}
	githubTagUrl += "/tags"
	// 每隔5分钟自动获取更新版本号信息
	gtimer.SetInterval(5*time.Minute, func() {
		updateVersion()
	})
}

// 获得最新GF框架版本号
func Latest(r *ghttp.Request) {
	if latestVersion.IsEmpty() {
		updateVersion()
	}
	r.Response.Write(latestVersion.String())
}

// 从github主仓库上解析获得最新GF框架版本号
func updateVersion() {
	content := ghttp.GetContent(githubTagUrl)
	array, _ := gregex.MatchAllString(`/gogf/gf/releases/tag/v([\.\d]+)"`, content)
	if len(array) > 0 && len(array[0]) == 2 {
		latestVersion.Set(array[0][1])
	}
}
