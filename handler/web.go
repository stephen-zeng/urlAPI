package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
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
	var newRequest request.Request
	web{}.requestBuilder(c, &newRequest)
	web{}.checker(&newRequest)
	if newRequest.Security.General.Unsafe {
		log.Printf("%s from %s\n", newRequest.Security.General.Info, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": newRequest.Security.General.Info,
		})
		return
	}
	beforeTask(&newRequest)
	webProcessor(&newRequest)
	afterTask(&newRequest)
	returner(c, &newRequest)
	return
}

func webProcessor(r *request.Request) {
	if r.DB.Task.Status == "success" {
		return
	}
	util.ErrorPrinter(r.Processor.WebImg.Process(&r.DB.Task))
	r.Processor.Return = r.Processor.WebImg.Return
}
