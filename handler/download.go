package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
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
		log.Printf("%s from %s\n", downloadRequest.Security.General.Info, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": downloadRequest.Security.General.Info,
		})
		return
	}
	util.ErrorPrinter(downloadRequest.Processor.Download.Process(&downloadRequest.DB.Task))
	taskSaver(&downloadRequest)
	downloadReturn(c, &downloadRequest)
	return
}

func downloadRequestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	ip := c.ClientIP()
	target := c.Query("img")
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	typ := util.TypeMap["download"]
	region := util.GetRegion(ip)
	r.Processor.Download = processor.Download{
		Target: target,
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
	r.Security.General.FrequencyChecker()
	r.Security.General.InfoChecker()
	r.Security.General.ExceptionChecker()
}

func downloadReturn(c *gin.Context, r *request.Request) {
	if r.Processor.Download.ReturnError != "" {
		c.Redirect(http.StatusFound, r.Processor.Download.ReturnError)
		return
	}
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="download.png"`))
	c.Header("Accept-Length", fmt.Sprintf("%d", len(r.Processor.Download.Return)))
	c.Writer.Write(r.Processor.Download.Return)
	return
}
