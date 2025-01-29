package router

import (
	"backend/cmd/txt"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func HttpServer() {
	r := gin.Default()
	r.GET("/txt", func(c *gin.Context) {
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		typ := c.Query("type")
		response, err := txt.Request(format, api, model, prompt, typ, "127.0.0.1", "www.goforit.top")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			if format == "json" {
				c.JSON(200, json.RawMessage(response))
			} else {
				c.String(200, response)
			}
		}
	})
	r.GET("/basket", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ball",
		})
	})
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(302, "https://www.google.com")
	})
	r.Run(":8080")
}
