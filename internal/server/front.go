package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Start(Port string) {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}
	r = gin.New()
	r.Use(cors.New(config))
	setAPI()
	setSession()
	setDash()
	r.Run(":" + Port)
}
