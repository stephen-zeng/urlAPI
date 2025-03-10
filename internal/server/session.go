package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"urlAPI/cmd/session"
)

func convertInterfaceToString(ori interface{}) [][]string {
	var ret [][]string
	if ori == nil {
		return ret
	}
	for _, i := range ori.([]interface{}) {
		var tmp []string
		for _, j := range i.([]interface{}) {
			if j == nil {
				tmp = append(tmp, "")
			} else {
				tmp = append(tmp, j.(string))
			}
		}
		ret = append(ret, tmp)
	}
	return ret
}

func setSession() {
	r.POST("/session", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		var dat sessionResp
		err := c.ShouldBindJSON(&dat)
		token := c.Request.Header.Get("Authorization")
		if err != nil { // auth err
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		typ, err := session.Auth(session.SessionConfig(
			session.WithSessionToken(token),
			session.WithSessionTime(time.Now()),
			session.WithSessionIP(c.ClientIP())))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		sessionConfig := session.Config{
			Operation:    dat.Operation,
			SessionToken: token,
			SessionType:  typ,
		}
		switch dat.Operation {
		case "login":
			sessionConfig.SessionTerm = dat.LoginTerm
		case "logout":
			sessionConfig.SessionTerm = dat.LoginTerm
		case "exit":
			sessionConfig.SessionTerm = dat.LoginTerm
		case "fetchTask":
			sessionConfig.TaskBy = dat.TaskBy
			sessionConfig.TaskCatagory = dat.TaskCatagory
		case "fetchSetting":
			sessionConfig.SettingPart = dat.SettingPart
		case "editSetting":
			sessionConfig.SettingPart = dat.SettingPart
			sessionConfig.SettingEdit = dat.SettingEdit
		case "newRepo":
			sessionConfig.RepoAPI = dat.RepoAPI
			sessionConfig.RepoInfo = dat.RepoInfo
		case "refreshRepo":
			sessionConfig.RepoUUID = dat.RepoUUID
		case "delRepo":
			sessionConfig.RepoUUID = dat.RepoUUID
		}

		response, err := session.New(sessionConfig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response.SessionIP = c.ClientIP()
		c.JSON(http.StatusOK, response)
	})
}
