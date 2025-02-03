package router

import (
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed template/*
var web embed.FS

func getContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".html"):
		return "text/html"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	default:
		return "text/plain"
	}
}

func dashboard() {
	subFS, _ := fs.Sub(web, "/")
	r.GET("/dash/*filepath", func(c *gin.Context) {
		filepath := strings.TrimPrefix(c.Param("filepath"), "/")

		// 处理根路径和默认文件
		if filepath == "" || strings.HasSuffix(filepath, "/") {
			filepath = "index.html"
		}

		// 尝试打开文件
		f, err := subFS.Open(filepath)
		if err != nil {
			// 处理前端路由：返回index.html
			if errors.Is(err, fs.ErrNotExist) && !strings.Contains(filepath, ".") {
				c.FileFromFS("index.html", http.FS(subFS))
				return
			}
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		defer f.Close()

		// 检查文件状态
		stat, err := f.Stat()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 禁止目录访问
		if stat.IsDir() {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 正确设置Content-Type
		contentType := getContentType(filepath)
		c.DataFromFS(http.StatusOK, contentType, f, http.FS(subFS))
	})
}
