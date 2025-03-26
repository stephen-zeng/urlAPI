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
		c.JSON(http.StatusForbidden, gin.H{
			"error": downloadRequest.Security.General.Info,
		})
		return
	}
	donwloadProcessor(&downloadRequest)
	downloadTaskSaver(&downloadRequest)
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
	r.Security.Operation = &r.Security.General
	r.Security.Operation.FrequencyChecker()
	r.Security.Operation.InfoChecker()
	r.Security.Operation.ExceptionChecker()
}

func donwloadProcessor(r *request.Request) {
	r.Processor.Operation = &r.Processor.Download
	err := r.Processor.Operation.Process(&r.DB.Task)
	if err != nil {
		log.Println(err)
	}
}

func downloadTaskSaver(r *request.Request) {
	if r.Security.General.SkipDB {
		return
	}
	r.DB.Operation = &r.DB.Task
	err := r.DB.Operation.Create()
	if err != nil {
		log.Println(err)
	}
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
