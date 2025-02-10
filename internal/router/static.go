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

func dashboard() {
	rootFS, _ := fs.Sub(webFS, "template/dist")
	assetsFS, _ := fs.Sub(webFS, "template/dist/assets")
	r.StaticFS("/assets", http.FS(assetsFS))
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	r.SetHTMLTemplate(tpl)
	r.GET("/dash", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
