package doc

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gopkg.in/russross/blackfriday.v2"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gregex"
    "fmt"
    "strings"
    "gitee.com/johng/gf/g/util/gstr"
)

// 获得指定uri路径的markdown文件内容
func GetMarkdown(path string) string {
    mdRoot  := g.Config().GetString("doc.path")
    content := gfile.GetContents(mdRoot + gfile.Separator + path + ".md")
    return content
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
    content    = string(blackfriday.Run([]byte(content)))
    pattern   := `href="(.+?)"`
    content, _ = gregex.ReplaceStringFunc(pattern, content, func(s string) string {
        match, _ := gregex.MatchString(pattern, gstr.Replace(s, ".md", ""))
        if len(match) > 1 {
            if match[1][0] != '/' && match[1][0] != '#' && !strings.Contains(match[1], "://") {
                return fmt.Sprintf(`href="/%s"`, match[1])
            }
        }
        return s
    })
    return content
}
