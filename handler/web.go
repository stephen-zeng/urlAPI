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
		c.JSON(http.StatusForbidden, gin.H{
			"error": webRequest.Security.General.Info,
		})
		return
	}
	if err := webRequest.Processor.Download.Process(&webRequest.DB.Task); err != nil {
		log.Println(err)
	}
	taskSaver(&webRequest)
	returner(c, webRequest.DB.Task.Return, webRequest.Processor.WebImg.Return)
	return
}

func webChecker(r *request.Request) {
	r.Security.General.FrequencyChecker()
	r.Security.General.InfoChecker()
	r.Security.General.ExceptionChecker()
	r.Security.WebImg.FunctionChecker(&r.Security.General)
	r.Security.WebImg.APIChecker(&r.Security.General)
}

func webBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	urlParse, _ := url.Parse(referer)
	host := getScheme(c) + urlParse.Host
	ip := c.ClientIP()
	target := c.Query("img")
	urlParse, _ = url.Parse(target)
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
