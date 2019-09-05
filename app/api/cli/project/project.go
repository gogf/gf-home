package project

import (
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"time"
)

var (
	projectUrl         = g.Config().GetString("project.url")
	projectMd5CacheKey = "project_zip_md5_" + projectUrl
	projectZipCacheKey = "project_zip_data_" + projectUrl
)

func init() {
	if projectUrl == "" {
		glog.Fatal("Project configuration cannot be empty")
	}
	// 每隔5分钟刷新项目zip缓存信息
	gtimer.SetInterval(5*time.Minute, func() {
		glog.Cat("project").Println("start refresh")
		refreshZipData()
		refreshMd5Data()
		glog.Cat("project").Println("start refresh")
	})
}

// 获得最新空项目zip文件md5值
func Md5(r *ghttp.Request) {
	md5 := gcache.GetOrSetFunc(projectMd5CacheKey, func() interface{} {
		data := getZipData()
		if len(data) == 0 {
			return nil
		}
		s, _ := gmd5.EncryptBytes(data)
		return s
	}, 0)
	r.Response.Write(md5)
}

// 获得最新空项目zip文件内容
func Zip(r *ghttp.Request) {
	r.Response.Write(getZipData())
}

// 下载并缓存空项目zip文件数据
func getZipData() []byte {
	data := gcache.GetOrSetFunc(projectZipCacheKey, func() interface{} {
		return ghttp.GetBytes(projectUrl)
	}, 0)
	return data.([]byte)
}

// 刷新项目zip缓存
func refreshMd5Data() {
	data := getZipData()
	if len(data) > 0 {
		s, _ := gmd5.EncryptBytes(data)
		gcache.Set(projectMd5CacheKey, s, 0)
	}
}

// 刷新项目zip缓存
func refreshZipData() {
	data := ghttp.GetBytes(projectUrl)
	if len(data) > 0 {
		gcache.Set(projectZipCacheKey, data, 0)
	}
}
