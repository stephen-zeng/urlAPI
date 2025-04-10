package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/url"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/security"
	"urlAPI/util"
)

func (txt) requestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
	ip := c.ClientIP()
	target := c.Query("prompt")
	model := c.Query("model")
	api := c.Query("api")
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["txt.gen"],
		Target:  target,
		Time:    time.Now(),
	}
	r.DB.Task = database.Task{
		UUID:    uuid.New().String(),
		Time:    time.Now(),
		IP:      ip,
		Type:    util.TypeMap["txt.gen"],
		Target:  target,
		Region:  region,
		Referer: referer,
		Device:  device,
		API:     api,
		Model:   model,
	}
	r.Processor.TxtGen = processor.TxtGen{
		API:    api,
		Model:  model,
		Target: target,
		Host:   host,
	}
	r.Processor.Filter = processor.TaskQueueFilter{
		Type:   "txt.gen",
		Target: target,
		API:    api,
	}
}

func (rand) requestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	ip := c.ClientIP()
	target := c.Query("user") + "/" + c.Query("repo")
	api := c.Query("api")
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["rand"],
		Target:  target,
		Time:    time.Now(),
	}
	r.Security.Rand = security.Rand{
		API:    api,
		Target: target,
	}
	r.DB.Task = database.Task{
		UUID:    uuid.New().String(),
		Time:    time.Now(),
		IP:      ip,
		Type:    util.TypeMap["rand"],
		Target:  target,
		Region:  region,
		Referer: referer,
		Device:  device,
		API:     api,
	}
	r.Processor.Rand = processor.Rand{
		API:    api,
		Target: target,
	}
}

func (img) requestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
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
		UUID:    uuid.New().String(),
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
	r.Processor.Filter = processor.TaskQueueFilter{
		Type:   "img.gen",
		Size:   size,
		Target: target,
		API:    api,
	}
}

func (web) requestBuilder(c *gin.Context, r *request.Request) {
	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
	ip := c.ClientIP()
	target := c.Query("img")
	urlParse, _ := url.Parse(target)
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
		UUID:    uuid.New().String(),
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
	r.Processor.Filter = processor.TaskQueueFilter{
		Type:   "web.img",
		Target: target,
		API:    api,
	}
}
