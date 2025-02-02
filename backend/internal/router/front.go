package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Start() {
	//gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}
	r = gin.New()
	r.Use(cors.New(config))
	setAPI()
	sessionListener()
	r.Run(":8080")
}
