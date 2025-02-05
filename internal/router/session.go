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
	r.POST("/dashsession", func(c *gin.Context) {
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
		var loginTerm bool
		var settingEdit [][]string
		var settingPart string
		var taskBy string
		var taskCatagory string
		var repoAPI string
		var repoInfo string
		var repoUUID string
		switch dat["operation"].(string) {
		case "login":
			loginTerm = dat["login_term"].(bool)
		case "logout":
			loginTerm = dat["login_term"].(bool)
		case "exit":
			loginTerm = dat["login_term"].(bool)
		case "fetchTask":
			taskCatagory = dat["task_catagory"].(string)
			taskBy = dat["task_by"].(string)
		case "fetchSetting":
			settingPart = dat["setting_part"].(string)
		case "editSetting":
			settingEdit = convertInterfaceToString(dat["setting_edit"])
			settingPart = dat["setting_part"].(string)
		case "newRepo":
			repoAPI = dat["repo_api"].(string)
			repoInfo = dat["repo_info"].(string)
		case "refreshRepo":
			repoUUID = dat["repo_uuid"].(string)
		case "delRepo":
			repoUUID = dat["repo_uuid"].(string)
		}

		response, err := session.New(session.SessionConfig(
			session.WithOperation(dat["operation"].(string)),
			session.WithSessionToken(token),
			session.WithSessionType(typ),
			session.WithSessionTerm(loginTerm),
			session.WithTaskCatagory(taskCatagory),
			session.WithTaskBy(taskBy),
			session.WithSettingPart(settingPart),
			session.WithSettingEdit(settingEdit),
			session.WithRepoAPI(repoAPI),
			session.WithRepoInfo(repoInfo),
			session.WithRepoUUID(repoUUID),
		))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		response.SessionIP = c.ClientIP()
		c.JSON(http.StatusOK, response)
	})
}
