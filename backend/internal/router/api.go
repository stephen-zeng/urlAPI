package router

import (
	"backend/cmd/img"
	"backend/cmd/txt"
	"backend/internal/file"
	"backend/internal/plugin"
	"backend/internal/security"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https://"
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http://"
}

func txtRequest() {
	r.GET("/txt", func(c *gin.Context) {
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		typ := c.Query("type")
		response, err := txt.Request(format, api, model, prompt, typ, c.ClientIP(), "www.goforit.top")
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
}

func imgRequest() {
	r.GET("/img", func(c *gin.Context) {
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		size := c.Query("size")
		response, err := img.Request(format, api, model, prompt, size,
			"127.0.0.1", "www.goforit.top", getScheme(c)+c.Request.Host)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			if format == "json" {
				c.JSON(200, json.RawMessage(response))
			} else {
				c.Redirect(302, response)
			}
		}
	})
}

func download() {
	r.GET("/download", func(c *gin.Context) {
		var err error
		err = security.NewRequest(security.SecurityConfig(
			security.WithIP(c.ClientIP()),
			security.WithDomain("www.goforit.top")))
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		Img := c.Query("img")
		Md := c.Query("md")
		var data []byte
		var suffix string
		if Img != "" {
			suffix = ".png"
			data, err = file.Fetch(file.FileConfig(
				file.WithType("img"),
				file.WithUUID(Img)))
			c.Header("Content-Type", "image/png")
		} else if Md != "" {
			suffix = ".md"
			data, err = file.Fetch(file.FileConfig(
				file.WithType("md"),
				file.WithUUID(Md)))
			c.Header("Content-Type", "text/plain")
		} else {
			c.JSON(400, gin.H{
				"error": "Unknow file type",
			})
		}
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="download%s"`, suffix))
			c.Header("Accept-Length", fmt.Sprintf("%d", len(data)))
			c.Writer.Write(data)
		}
	})
}

func rand() {
	r.GET("/rand", func(c *gin.Context) {
		err := security.NewRequest(security.SecurityConfig(
			security.WithIP(c.ClientIP()),
			security.WithDomain("www.goforit.top")))
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		API := c.Query("api")
		User := c.Query("user")
		Repo := c.Query("repo")
		response, err := plugin.Request(plugin.PluginConfig(
			plugin.WithAPI(API),
			plugin.WithUser(User),
			plugin.WithRepo(Repo)))
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		c.Redirect(302, response)
	})
}

func setAPI() {
	txtRequest()
	imgRequest()
	download()
	rand()
}
