package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"urlAPI/request"
	"urlAPI/util"
)

func randHandler(c *gin.Context) {
	var newRequest request.Request
	rand{}.requestBuilder(c, &newRequest)
	rand{}.checker(&newRequest)
	if newRequest.Security.General.Unsafe {
		log.Printf("%s from %s\n", newRequest.Security.General.Info, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": newRequest.Security.General.Info,
		})
		return
	}
	util.ErrorPrinter(newRequest.Processor.Rand.Process(&newRequest.DB.Task))
	newRequest.Processor.Return = newRequest.Processor.Rand.Return
	if !newRequest.Security.General.SkipDB {
		util.ErrorPrinter(newRequest.DB.Task.Create())
	}
	returner(c, &newRequest)
}
