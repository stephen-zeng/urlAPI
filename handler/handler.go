package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func Handler(port string) {
	gin.SetMode(gin.ReleaseMode)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST"}
	r := gin.New()
	r.Use(cors.New(config))
	r.GET("/txt", txtHandler)
	r.GET("/img", imgHandler)
	r.GET("/rand", randHandler)
	r.GET("/web", webHandler)
	r.GET("/download", downloadHandler)
	r.POST("/session", sessionHandler)
	log.Printf("The server will be running on port %s", port)
	r.Run(":" + port)
}
