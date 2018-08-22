package doc

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gopkg.in/russross/blackfriday.v2"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gstr"
)

// 获得指定uri路径的markdown文件内容
func GetMarkdown(path string) string {
    mdRoot  := g.Config().GetString("gf-doc.path")
    content := gfile.GetContents(mdRoot + gfile.Separator + path + ".md")
    return content
}

// 获得解析为html的markdown文件内容
func GetParsed(path string) string {
    return ParseMarkdown(GetMarkdown(path))
}

// 解析markdown为html
func ParseMarkdown(content string) string {
    content = string(blackfriday.Run([]byte(content)))
    content = gstr.ReplaceByMap(content, map[string]string{
        `href="` : `href="/`,
    })
    return content
}
