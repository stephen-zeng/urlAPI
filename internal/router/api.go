package router

import (
	"backend/cmd/img"
	"backend/cmd/txt"
	"backend/internal/plugin"
	"backend/internal/security"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
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
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		domain := referer.Hostname()
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		typ := c.Query("type")
		response, err := txt.Request(format, api, model, prompt, typ, c.ClientIP(), domain)
		if err != nil {
			log.Println(err)
			c.String(200, err.Error())
		} else {
			if format == "json" {
				c.JSON(200, response)
			} else {
				c.String(200, response.Response)
			}
		}
	})
}

func imgRequest() {
	r.GET("/img", func(c *gin.Context) {
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		domain := referer.Hostname()
		format := c.Query("format")
		api := c.Query("api")
		model := c.Query("model")
		prompt := c.Query("prompt")
		size := c.Query("size")
		response, err := img.GenRequest(
			c.ClientIP(), domain,
			model, api, prompt, size,
			getScheme(c)+c.Request.Host)
		if err != nil {
			log.Println(err)
			c.Redirect(302, response.URL)
		} else {
			if format == "json" {
				c.JSON(200, response)
			} else {
				c.Redirect(302, response.URL)
			}
		}
	})
}

func rand() {
	r.GET("/rand", func(c *gin.Context) {
		referer, err := url.Parse(c.Request.Referer())
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		domain := referer.Hostname()
		err = security.NewRequest(security.SecurityConfig(
			security.WithIP(c.ClientIP()),
			security.WithDomain(domain)))
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
		c.Redirect(302, response.URL)
	})
}

func setAPI() {
	txtRequest()
	imgRequest()
	rand()
}
