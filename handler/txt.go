package handler

import (
	"github.com/gin-gonic/gin"
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

func txtHandler(c *gin.Context) {
	var txtRequest request.Request
	txtRequestBuilder(c, &txtRequest)
	txtChecker(&txtRequest)
	if txtRequest.Security.General.Unsafe {
		c.JSON(http.StatusForbidden, gin.H{
			"error": txtRequest.Security.General.Info,
		})
		return
	}
	if txtOldTask(&txtRequest) {
		c.Redirect(http.StatusFound, txtRequest.Processor.TxtGen.Return)
		return
	}
	if err := txtRequest.Processor.Rand.Process(&txtRequest.DB.Task); err != nil {
		log.Println(err)
	}
	taskSaver(&txtRequest)
	returner(c, txtRequest.DB.Task.Return, txtRequest.Processor.Rand.Return)
	return
}

func txtOldTask(r *request.Request) bool {
	var hasOldTask bool
	expireTime, _ := strconv.Atoi(database.SettingMap["txt"][3])
	taskFinder := database.Task{
		Type:   r.DB.Task.Type,
		Status: "success",
	}
	taskDBList, err := taskFinder.Read()
	if err != nil {
		log.Printf("Handler txtHandler %s", err.Error())
		return false
	}
	for _, task := range taskDBList.TaskList {
		os.Remove(processor.ImgPath + task.UUID + ".png")
		if time.Now().Sub(task.Time) <= time.Duration(expireTime)*time.Minute {
			hasOldTask = true
			r.Processor.TxtGen.Return = r.Processor.TxtGen.Host + "download?img=" + task.UUID
		}
	}
	return hasOldTask
}

func txtChecker(r *request.Request) {
	r.Security.General.FrequencyChecker()
	r.Security.General.ExceptionChecker()
	r.Security.General.InfoChecker()
	r.Security.TxtGen.APIChecker(&r.Security.General)
	r.Security.TxtGen.FunctionChecker(&r.Security.General)
}

func txtRequestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	urlParse, _ := url.Parse(referer)
	host := getScheme(c) + urlParse.Host
	ip := c.ClientIP()
	target := c.Query("prompt")
	model := c.Query("model")
	api := c.Query("api")
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["txt.gen"],
		Target:  target,
		Time:    time.Now(),
	}
	r.Security.TxtSum = security.TxtSum{
		Model: model,
		API:   api,
	}
	r.DB.Task = database.Task{
		Time:    time.Now(),
		IP:      ip,
		Type:    util.TypeMap["txt.gen"],
		Target:  target,
		Region:  region,
		Referer: referer,
		Device:  device,
		API:     api,
		Model:   model,
	}
	r.Processor.TxtGen = processor.TxtGen{
		API:    api,
		Model:  model,
		Target: target,
		Host:   host,
	}
}
