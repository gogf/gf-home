package document

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gfcache"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/russross/blackfriday/v2"
)

var (
	// Path is the document diractory for gf-doc.
	Path = g.Config().GetString("document.path")
)

var (
	// Cache for documents.
	cache = gcache.New()
)

func init() {
	// Path checking..
	if Path == "" {
		glog.Fatal("configuration for document.path cannot be empty")
	}
}

// UpdateDocGit does the "git pull" for gf-doc.
func UpdateDocGit() {
	err := gproc.ShellRun(
		fmt.Sprintf(`cd %s && git pull origin master`, Path),
	)
	if err == nil {
		cache.Clear()

		glog.Cat("doc-hook").Printf("doc hook updates")
	} else {
		glog.Cat("doc-hook").Printf("doc hook updates error: %v", err)
	}
}

// SearchMdByKey searches the markdown files with specified key
// and returns the file list.
func SearchMdByKey(key string) []string {
	glog.Cat("search").Println(key)
	v := cache.GetOrSetFunc("doc_search_result_"+key, func() interface{} {
		// 当该key的检索缓存不存在时，执行检索
		var (
			array   = garray.NewStrArray(true)
			docPath = g.Config().GetString("document.path")
			paths   = cache.GetOrSetFunc("doc_files_recursive", func() interface{} {
				// 当目录列表不存在时，执行检索
				paths, _ := gfile.ScanDir(docPath, "*.md", true)
				return paths
			}, 0)
		)
		// 遍历markdown文件列表，执行字符串搜索
		for _, path := range gconv.Strings(paths) {
			content := gfcache.GetContents(path)
			if len(content) > 0 {
				if strings.Index(content, key) != -1 {
					index := gstr.Replace(path, ".md", "")

					// 能同时处理绝对路径相对路径
					index = index[strings.Index(index, docPath)+len(docPath):]

					// 只能处理绝对路径
					//index = gstr.Replace(index, docPath, "")

					// 替换斜杠，修复windows下搜索失败的问题
					index = strings.ReplaceAll(index, "\\", "/")
					array.Append(index)
				}
			}
		}
		return array.Slice()
	}, 0)

	return gconv.Strings(v)
}

// 根据path参数获得层级显示的title
func GetTitleByPath(path string) string {
	v := cache.GetOrSetFunc("title_by_path_"+path, func() interface{} {
		type lineItem struct {
			indent int
			name   string
		}
		path = strings.TrimLeft(path, "/")
		array := make([]lineItem, 0)
		mdContent := GetMarkdown("menus")
		lines := strings.Split(mdContent, "\n")
		indent := 0
		for _, line := range lines {
			match, _ := gregex.MatchString(`(\s*)\*\s+\[(.+)\]\((.+)\)`, line)
			if len(match) == 4 {
				item := lineItem{
					indent: len(match[1]),
					name:   match[2],
				}
				mdPath := gstr.Replace(match[3], ".md", "")
				if item.indent > indent || len(array) == 0 {
					array = append(array, item)
				} else if len(match[1]) == indent {
					array[len(array)-1] = item
				} else {
					newArray := make([]lineItem, 0)
					for _, v := range array {
						if v.indent < item.indent {
							newArray = append(newArray, v)
						}
					}
					newArray = append(newArray, item)
					array = newArray
				}
				indent = item.indent
				if mdPath == path {
					break
				}
			}
		}
		if len(array) > 0 {
			title := ""
			for i := len(array) - 1; i >= 0; i-- {
				if len(title) > 0 {
					title += " - " + array[i].name
				} else {
					title = array[i].name
				}
			}
			return title
		}
		return nil
	}, 0)
	if v != nil {
		return v.(string)
	}
	return ""
}

// 获得指定uri路径的markdown文件内容
func GetMarkdown(path string) string {
	return gfcache.GetContents(Path + gfile.Separator + path + ".md")
}

// 获得解析为html的markdown文件内容
func GetParsed(path string) string {
	return ParseMarkdown(GetMarkdown(path))
}

// 解析markdown为html
func ParseMarkdown(content string) string {
	if content == "" {
		return ""
	}
	// src及href 替换为/xxx模式的绝对连接
	content = string(blackfriday.Run([]byte(content)))
	pattern := `(src|href)=["'](.+?)["']`
	content, _ = gregex.ReplaceStringFunc(pattern, content, func(s string) string {
		match, _ := gregex.MatchString(pattern, gstr.Replace(s, ".md", ""))
		if len(match) > 1 {
			if match[2][0] != '/' && match[2][0] != '#' && !strings.Contains(match[2], "://") {
				return fmt.Sprintf(`%s="/%s"`, match[1], match[2])
			}
		}
		return s
	})
	return content
}
