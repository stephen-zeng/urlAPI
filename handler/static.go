package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func staticHandler(c *gin.Context) {
	switch {
	case c.Request.URL.Path == "/dash":
		c.HTML(http.StatusOK, "index.html", nil)
	case len(c.Request.URL.Path) > 5 && c.Request.URL.Path[:6] == "/dash/":
		c.HTML(http.StatusOK, "index.html", nil)
	default:
		c.Redirect(301, "https://www.bilibili.com/video/BV1GJ411x7h7/")
	}
}
