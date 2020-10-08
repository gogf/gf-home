package binary

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"time"
)

var (
	cliRoot  = g.Config().GetString("cli.path")
	cacheMap = gmap.NewStrStrMap(true)
	filesSet = gset.NewStrSet(true)
)

func init() {
	if cliRoot == "" {
		glog.Fatal("CLI configuration cannot be empty")
	}
	refreshFilesSet()
	gtimer.SetInterval(5*time.Minute, func() {
		refreshFilesSet()
		refreshCacheMap()
	})
}

// CLI二进制文件浏览
func Index(r *ghttp.Request) {
	path := cliRoot + "/" + r.GetString("path")
	if gfile.IsFile(path) {
		// 引导到CDN地址下载
		if cdnUrl := g.Config().GetString("cdn.url"); cdnUrl != "" {
			r.Response.RedirectTo(
				fmt.Sprintf(
					`%s%s?%s`,
					cdnUrl,
					r.URL.Path,
					cacheMap.Get(gfile.RealPath(path)),
				),
			)
		}
	}
	r.Response.ServeFile(cliRoot+"/"+path, true)
}

// 获得最新CLI工具二进制文件md5值
func Md5(r *ghttp.Request) {
	path := buildBinaryPath(r)
	if !filesSet.Contains(path) {
		return
	}
	md5 := cacheMap.GetOrSetFunc(path, func() string {
		s, _ := gmd5.EncryptFile(path)
		return s
	})
	r.Response.Write(md5)
}

// 根据请求参数
func buildBinaryPath(r *ghttp.Request) string {
	var (
		os   = r.Get("os")
		arch = r.Get("arch")
		path = r.GetString("path")
		name = "gf"
	)
	if os == "windows" {
		name += ".exe"
	}
	return gfile.Abs(cliRoot + "/" + path + fmt.Sprintf(`%s_%s/%s`, os, arch, name))
}

// 刷新文件md5缓存
func refreshCacheMap() {
	for _, path := range filesSet.Slice() {
		s, err := gmd5.EncryptFile(path)
		if err == nil {
			cacheMap.Set(path, s)
		} else {
			glog.Error(err)
		}
	}
}

// 刷新文件列表缓存
func refreshFilesSet() {
	files, err := gfile.ScanDirFile(cliRoot, "*", true)
	if err != nil {
		glog.Error(err)
	} else {
		filesSet.Clear()
		filesSet.Add(files...)
	}
}
