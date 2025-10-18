package controllers

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// RenderMarkdownFile 文档解析
func RenderMarkdownFile(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.HTML(http.StatusOK, "result.html", gin.H{
			"title": "404 - 文档不存在",
			"desc":  "检查路径是否正确或者返回首页",
		})
		return
	}
	// 读取文件内容
	md, err := ioutil.ReadFile("./docs/" + name + ".md")
	if err != nil {
		c.HTML(http.StatusOK, "result.html", gin.H{
			"title": "404 - 文档不存在",
			"desc":  "检查路径是否正确或者返回首页",
		})
		return
	}

	// 渲染 Markdown 为 HTML
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)
	var buf bytes.Buffer
	if err := markdown.Convert(md, &buf); err != nil {
		c.HTML(http.StatusOK, "result.html", gin.H{
			"title": "Error - 文档解析错误",
			"desc":  "反馈给我们或者返回首页",
		})
		return
	}
	content := buf.String()
	c.HTML(http.StatusOK, "page", gin.H{
		"title":   name,
		"content": template.HTML(content),
	})
}
