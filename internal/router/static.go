package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed template/dist/*
var webFS embed.FS

func setDash() {
	rootFS, _ := fs.Sub(webFS, "template/dist")
	assetsFS, _ := fs.Sub(webFS, "template/dist/assets")
	r.StaticFS("/assets", http.FS(assetsFS))
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	r.SetHTMLTemplate(tpl)
	r.NoRoute(func(c *gin.Context) {
		switch {
		case c.Request.URL.Path == "/dash":
			c.HTML(http.StatusOK, "index.html", nil)
		case len(c.Request.URL.Path) > 5 && c.Request.URL.Path[:6] == "/dash/":
			c.HTML(http.StatusOK, "index.html", nil)
		default:
			c.Redirect(301, "https://www.bilibili.com/video/BV1GJ411x7h7/")
		}
	})
}
