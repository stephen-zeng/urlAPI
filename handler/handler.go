package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"urlAPI/static"
)

var r *gin.Engine

func Handler(port string) {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}
	r = gin.New()
	r.Use(cors.New(config))
	rootFS, _ := fs.Sub(static.StaticFS, "dist")
	assetsFS, _ := fs.Sub(static.StaticFS, "dist/assets")
	r.StaticFS("/assets", http.FS(assetsFS))
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	r.SetHTMLTemplate(tpl)
	r.GET("/txt", txtHandler)
	r.GET("/img", imgHandler)
	r.GET("/rand", randHandler)
	r.GET("/web", webHandler)
	r.GET("/download", downloadHandler)
	r.POST("/session", sessionHandler)
	r.NoRoute(staticHandler)
	log.Printf("The server will be running on port %s", port)
	r.Run(":" + port)
}
