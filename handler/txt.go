package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
	"urlAPI/util"
)

func txtHandler(c *gin.Context) {
	var newRequest request.Request
	txt{}.requestBuilder(c, &newRequest)
	txt{}.checker(&newRequest)
	if newRequest.Security.General.Unsafe {
		log.Printf("%s from %s\n", newRequest.Security.General.Info, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": newRequest.Security.General.Info,
		})
		return
	}
	beforeTask(&newRequest)
	txtProcessor(&newRequest)
	afterTask(&newRequest)
	returner(c, &newRequest)
	return
}

func txtProcessor(r *request.Request) {
	if r.DB.Task.Status == "success" {
		return
	}
	util.ErrorPrinter(r.Processor.TxtGen.Process(&r.DB.Task))
	r.Processor.Return = r.Processor.TxtGen.Return
}
