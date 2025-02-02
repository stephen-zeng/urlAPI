package router

import (
	"backend/cmd/session"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func sessionListener() {
	r.POST("/session", func(c *gin.Context) {
		dat := make(map[string]interface{})
		err := c.ShouldBindJSON(&dat)
		token := c.Request.Header.Get("Authorization")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		typ, err := session.Auth(session.SessionConfig(session.WithToken(token)))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		fmt.Println(typ)
	})
}
