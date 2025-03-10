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
		dat := make(map[string]interface{})
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
			Operation:    dat["operation"].(string),
			SessionToken: token,
			SessionType:  typ,
		}
		switch dat["operation"].(string) {
		case "login":
			sessionConfig.SessionTerm = dat["login_term"].(bool)
		case "logout":
			sessionConfig.SessionTerm = dat["login_term"].(bool)
		case "exit":
			sessionConfig.SessionTerm = dat["login_term"].(bool)
		case "fetchTask":
			sessionConfig.TaskBy = dat["task_by"].(string)
			sessionConfig.TaskCatagory = dat["task_catagory"].(string)
		case "fetchSetting":
			sessionConfig.SettingPart = dat["setting_part"].(string)
		case "editSetting":
			sessionConfig.SettingPart = dat["setting_part"].(string)
			sessionConfig.SettingEdit = convertInterfaceToString(dat["setting_edit"])
		case "newRepo":
			sessionConfig.RepoAPI = dat["repo_api"].(string)
			sessionConfig.RepoInfo = dat["repo_info"].(string)
		case "refreshRepo":
			sessionConfig.RepoUUID = dat["repo_uuid"].(string)
		case "delRepo":
			sessionConfig.RepoUUID = dat["repo_uuid"].(string)
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
