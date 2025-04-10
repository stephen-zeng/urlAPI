package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
	"urlAPI/util"
)

func imgHandler(c *gin.Context) {
	var newRequest request.Request
	img{}.requestBuilder(c, &newRequest)
	img{}.checker(&newRequest)
	if newRequest.Security.General.Unsafe {
		log.Printf("%s from %s\n", newRequest.Security.General.Info, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": newRequest.Security.General.Info,
		})
		return
	}
	beforeTask(&newRequest)
	imgProcessor(&newRequest)
	afterTask(&newRequest)
	returner(c, &newRequest)
	return
}

func imgProcessor(r *request.Request) {
	if r.DB.Task.Status == "success" {
		return
	}
	util.ErrorPrinter(r.Processor.ImgGen.Process(&r.DB.Task))
	r.Processor.Return = r.Processor.ImgGen.Return
}
