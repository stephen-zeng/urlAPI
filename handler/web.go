package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/security"
	"urlAPI/util"
)

var webAPIMap = map[string]string{
	"github.com":       "github",
	"gitee.com":        "gitee",
	"www.bilibili.com": "bilibili",
	"www.youtube.com":  "youtube",
	"arxiv.org":        "arxiv",
	"www.ithome.com":   "ithome",
}

func webHandler(c *gin.Context) {
	var webRequest request.Request
	webBuilder(c, &webRequest)
	webChecker(&webRequest)
	if webRequest.Security.General.Unsafe {
		log.Println(webRequest.Security.General.Info)
		c.JSON(http.StatusForbidden, gin.H{
			"error": webRequest.Security.General.Info,
		})
		return
	}
	if webOldTask(&webRequest) {
		returner(c, webRequest.DB.Task.Return, webRequest.Processor.WebImg.Return)
		return
	}
	if err := webRequest.Processor.WebImg.Process(&webRequest.DB.Task); err != nil {
		log.Println(err)
	}
	taskSaver(&webRequest)
	returner(c, webRequest.DB.Task.Return, webRequest.Processor.WebImg.Return)
	return
}

func webOldTask(r *request.Request) bool {
	var hasOldTask bool
	expireTime, _ := strconv.Atoi(database.SettingMap["web"][3])
	taskFinder := database.Task{
		Target: r.DB.Task.Target,
		Type:   r.DB.Task.Type,
		Status: "success",
	}
	taskDBList, err := taskFinder.Read()
	if err != nil {
		log.Printf("Handler webHandler %s", err.Error())
		return false
	}
	for _, task := range taskDBList.TaskList {
		if time.Now().Sub(task.Time) <= time.Duration(expireTime)*time.Minute {
			hasOldTask = true
			r.Processor.WebImg.Return = r.Processor.WebImg.Host + "/download?img=" + task.UUID
			r.DB.Task.Return = task.Return
		} else {
			os.Remove(processor.ImgPath + task.UUID + ".png")
		}
	}
	return hasOldTask
}

func webChecker(r *request.Request) {
	r.Security.General.FrequencyChecker()
	r.Security.General.InfoChecker()
	r.Security.General.ExceptionChecker()
	r.Security.WebImg.FunctionChecker(&r.Security.General)
}

func webBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
	ip := c.ClientIP()
	target := c.Query("img")
	urlParse, _ := url.Parse(target)
	api := urlParse.Host
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["web.img"],
		Target:  target,
		Time:    time.Now(),
	}
	r.Security.WebImg = security.WebImg{
		API: api,
	}
	r.DB.Task = database.Task{
		UUID:    uuid.New().String(),
		Time:    time.Now(),
		IP:      ip,
		Type:    util.TypeMap["web.img"],
		Target:  target,
		Region:  region,
		Referer: referer,
		Device:  device,
		API:     api,
	}
	r.Processor.WebImg = processor.WebImg{
		API:    api,
		Target: target,
		Host:   host,
	}
}
