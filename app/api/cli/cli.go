package cli

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/crypto/gmd5"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gcache"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/glog"
	"time"
)

var (
	cliRoot = g.Config().GetString("cli.path")
)

func init() {
	if cliRoot == "" {
		glog.Fatal("CLI configuration cannot be empty")
	}
}

// CLI二进制文件浏览
func Index(r *ghttp.Request) {
	r.Response.ServeFile(cliRoot+"/"+r.Get("path"), true)
}

// 检查是否需要更新工具。
// 输出1: 需要更新, 0: 不更新, 其他：错误。
func Check(r *ghttp.Request) {
	inputMd5 := r.Get("md5")
	if inputMd5 == "" {
		r.Response.Write("md5 value cannot be empty")
		return
	}
	path := getFilePath(r)
	md5 := gcache.GetOrSetFuncLock("cli_binary_md5_"+path, func() interface{} {
		s, err := gmd5.EncryptFile(path)
		if err != nil {
			return nil
		}
		return s
	}, time.Hour)
	if md5 == nil {
		r.Response.Write("invalid params")
		return
	}
	if inputMd5 == md5 {
		r.Response.Write("0")
	} else {
		r.Response.Write("1")
	}
}

// CLI二进制文件下载
func Download(r *ghttp.Request) {
	glog.Cat("cli/download").Println(r.GetMap())

	path := getFilePath(r)
	r.Response.ServeFileDownload(path, gfile.Basename(path))
}

func getFilePath(r *ghttp.Request) string {
	os := r.Get("os")
	arch := r.Get("arch")
	name := "gf"
	if os == "windows" {
		name += ".exe"
	}
	return cliRoot + "/" + r.Get("path") + fmt.Sprintf(`%s_%s/%s`, os, arch, name)
}
