package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/security"
	"urlAPI/util"
)

func imgHandler(c *gin.Context) {
	var imgRequest request.Request
	imgBuilder(c, &imgRequest)
	imgChecker(&imgRequest)
	if imgRequest.Security.General.Unsafe {
		c.JSON(http.StatusForbidden, gin.H{
			"error": imgRequest.Security.General.Info,
		})
		return
	}
	if err := imgRequest.Processor.Download.Process(&imgRequest.DB.Task); err != nil {
		log.Println(err)
	}
	taskSaver(&imgRequest)
	returner(c, imgRequest.DB.Task.Return, imgRequest.Processor.ImgGen.Return)
	return
}

func imgChecker(r *request.Request) {
	r.Security.General.ExceptionChecker()
	r.Security.General.InfoChecker()
	r.Security.General.FrequencyChecker()
	r.Security.ImgGen.APIChecker(&r.Security.General)
	r.Security.ImgGen.FunctionChecker(&r.Security.General)
}

func imgBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	urlParse, _ := url.Parse(referer)
	host := urlParse.Host
	ip := c.ClientIP()
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	model := c.Query("model")
	target := c.Query("prompt")
	size := c.Query("size")
	api := c.Query("api")

	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["img.gen"],
		Target:  target,
		Time:    time.Now(),
	}
	r.Security.ImgGen = security.ImgGen{
		API:   api,
		Model: model,
	}
	r.DB.Task = database.Task{
		Time:    time.Now(),
		IP:      ip,
		Type:    util.TypeMap["img.gen"],
		Target:  target,
		Region:  region,
		Referer: referer,
		Device:  device,
		API:     api,
		Model:   model,
		Size:    size,
	}
	r.Processor.ImgGen = processor.ImgGen{
		API:    api,
		Model:  model,
		Target: target,
		Host:   host,
		Size:   size,
	}
}
