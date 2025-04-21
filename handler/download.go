package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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
	afterTask(&downloadRequest)
	downloadReturn(c, &downloadRequest)
	return
}

func downloadRequestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	ip := c.ClientIP()
	target := c.Query("img")
	typ := util.TypeMap["download"]
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
}

func downloadChecker(r *request.Request) {
	r.Security.General.SkipDB = true
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
