package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/security"
	"urlAPI/util"
)

func randHandler(c *gin.Context) {
	var randRequest request.Request
	randRequestBuilder(c, &randRequest)
	randChecker(&randRequest)
	if randRequest.Security.General.Unsafe {
		c.JSON(http.StatusForbidden, gin.H{
			"error": randRequest.Security.General.Info,
		})
		return
	}
	if err := randRequest.Processor.Rand.Process(&randRequest.DB.Task); err != nil {
		log.Println(err)
	}
	taskSaver(&randRequest)
	returner(c, randRequest.DB.Task.Return, randRequest.Processor.Rand.Return)
}

func randChecker(r *request.Request) {
	r.Security.General.FrequencyChecker()
	r.Security.General.InfoChecker()
	r.Security.General.ExceptionChecker()
	r.Security.Rand.FunctionChecker(&r.Security.General)
	r.Security.Rand.APIChecker(&r.Security.General)
}

func randRequestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	ip := c.ClientIP()
	target := c.Query("user") + "/" + c.Query("repo")
	api := c.Query("api")
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["rand"],
		Target:  target,
		Time:    time.Now(),
	}
	r.Security.Rand = security.Rand{
		API:    api,
		Target: target,
	}
	r.DB.Task = database.Task{
		Time:    time.Now(),
		IP:      ip,
		Type:    util.TypeMap["rand"],
		Target:  target,
		Region:  region,
		Referer: referer,
		Device:  device,
		API:     api,
	}
	r.Processor.Rand = processor.Rand{
		API:    api,
		Target: target,
	}
}
