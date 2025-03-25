package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/security"
	"urlAPI/util"
)

func downloadHandler(c *gin.Context) {
	var downloadRequest request.Request
	downloadRequestBuilder(c, &downloadRequest)
	downloadChecker(&downloadRequest)
	if downloadRequest.Security.General.Unsafe {
		c.JSON(http.StatusForbidden, gin.H{
			"error": downloadRequest.Security.General.Info,
		})
	}
}

func downloadRequestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	ip := c.ClientIP()
	target := c.Query("img")
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	typ := util.TypeMap["download"]
	region := util.GetRegion(ip)
	r.Processor.API = processor.API{
		Referer: referer,
		IP:      ip,
		Target:  target,
		Type:    typ,
		Device:  device,
	}
	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Target:  target,
		Type:    typ,
		Time:    time.Now(),
	}
	r.DB.Task = database.Task{
		UUID:   uuid.New().String(),
		Time:   time.Now(),
		IP:     ip,
		Target: target,
		Type:   typ,
		Device: device,
		Region: region,
	}
}

func downloadChecker(r *request.Request) {
	r.Security.Operation = &r.Security.General
	r.Security.Operation.FrequencyChecker()
	r.Security.Operation.RefererChecker()
	r.Security.Operation.ExceptionChecker()
}
