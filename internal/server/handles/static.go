package handles

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticHandler(c *gin.Context) {
	switch {
	case c.Request.URL.Path == "/v1" || strings.HasPrefix(c.Request.URL.Path, "/v1/"):
		log.Printf("[NoRoute] Unsupported API route: %s %s | IP: %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"message": "unsupported API route: " + c.Request.URL.Path,
				"type":    "invalid_request_error",
				"code":    "unsupported_route",
			},
		})
	case c.Request.URL.Path == "/dash":
		c.HTML(http.StatusOK, "index.html", nil)
	case len(c.Request.URL.Path) > 5 && c.Request.URL.Path[:6] == "/dash/":
		c.HTML(http.StatusOK, "index.html", nil)
	default:
		c.Redirect(301, "https://www.bilibili.com/video/BV1GJ411x7h7/")
	}
}

func MethodNotAllowedHandler(c *gin.Context) {
	log.Printf("[NoMethod] Unsupported method for route: %s %s | IP: %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": gin.H{
			"message": "method not allowed: " + c.Request.Method + " " + c.Request.URL.Path,
			"type":    "invalid_request_error",
			"code":    "method_not_allowed",
		},
	})
}
