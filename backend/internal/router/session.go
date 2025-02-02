package router

import (
	"backend/cmd/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func convertInterfaceToString(ori interface{}) [][]string {
	var ret [][]string
	if ori == nil {
		return ret
	}
	for _, i := range ori.([]interface{}) {
		var tmp []string
		for _, j := range i.([]interface{}) {
			tmp = append(tmp, j.(string))
		}
		ret = append(ret, tmp)
	}
	return ret
}

func sessionListener() {
	r.POST("/session", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		dat := make(map[string]interface{})
		err := c.ShouldBindJSON(&dat)
		token := c.Request.Header.Get("Authorization")
		if err != nil { // auth err
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		typ, err := session.Auth(session.SessionConfig(
			session.WithToken(token),
			session.WithTime(time.Now()),
			session.WithIP(c.ClientIP())))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var edit [][]string
		var part string
		var term bool
		task := ""
		if _, ok := dat["edit"]; ok {
			edit = convertInterfaceToString(dat["edit"])
		}
		if _, ok := dat["part"]; ok {
			part = dat["part"].(string)
		}
		if _, ok := dat["term"]; ok {
			term = dat["term"].(bool)
		}
		if _, ok := dat["task"]; ok {
			task = dat["task"].(string)
		}
		response, err := session.New(session.SessionConfig(
			session.WithToken(token),
			session.WithType(typ),
			session.WithTerm(term),
			session.WithOperation(dat["operation"].(string)),
			session.WithPart(part),
			session.WithEdit(edit),
			session.WithBy(task)))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})
}
