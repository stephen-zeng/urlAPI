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
	var query apiQuery
	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
	ip := c.ClientIP()
	_ = c.ShouldBindQuery(&query)
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Info:    query.More,
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["txt.gen"],
		Target:  query.Prompt,
		Time:    time.Now(),
	}
	r.DB.Task = database.Task{
		UUID:     uuid.New().String(),
		Time:     time.Now(),
		IP:       ip,
		Type:     util.TypeMap["txt.gen"],
		Target:   query.Prompt,
		Region:   region,
		Referer:  referer,
		Device:   device,
		API:      query.API,
		Model:    query.Model,
		MoreInfo: query.More,
	}
	r.Processor.TxtGen = processor.TxtGen{
		API:    query.API,
		Model:  query.Model,
		Target: query.Prompt,
		Host:   host,
	}
	r.Processor.Filter = processor.TaskQueueFilter{
		Type:   "txt.gen",
		Target: query.Prompt,
		API:    query.API,
	}
}

func (rand) requestBuilder(c *gin.Context, r *request.Request) {
	var query apiQuery
	_ = c.ShouldBindQuery(&query)
	referer := c.Request.Referer()
	ip := c.ClientIP()
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)
	r.Security.General = security.General{
		Info:    query.More,
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["rand"],
		Target:  query.User + "/" + query.Repo,
		Time:    time.Now(),
	}
	r.Security.Rand = security.Rand{
		API:    query.API,
		Target: query.Prompt,
	}
	r.DB.Task = database.Task{
		UUID:     uuid.New().String(),
		Time:     time.Now(),
		IP:       ip,
		Type:     util.TypeMap["rand"],
		Target:   query.User + "/" + query.Repo,
		Region:   region,
		Referer:  referer,
		Device:   device,
		API:      query.API,
		MoreInfo: query.More,
	}
	r.Processor.Rand = processor.Rand{
		API:    query.API,
		Target: query.User + "/" + query.Repo,
	}
}

func (img) requestBuilder(c *gin.Context, r *request.Request) {
	var query apiQuery
	_ = c.ShouldBindQuery(&query)
	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
	ip := c.ClientIP()
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)

	r.Security.General = security.General{
		Info:    query.More,
		Referer: referer,
		IP:      ip,
		Type:    util.TypeMap["img.gen"],
		Target:  query.Prompt,
		Time:    time.Now(),
	}
	r.Security.ImgGen = security.ImgGen{
		API:   query.API,
		Model: query.Model,
	}
	r.DB.Task = database.Task{
		UUID:     uuid.New().String(),
		Time:     time.Now(),
		IP:       ip,
		Type:     util.TypeMap["img.gen"],
		Target:   query.Prompt,
		Region:   region,
		Referer:  referer,
		Device:   device,
		API:      query.API,
		Model:    query.Model,
		MoreInfo: query.More,
		Size:     query.Size,
	}
	r.Processor.ImgGen = processor.ImgGen{
		API:    query.API,
		Model:  query.Model,
		Target: query.Prompt,
		Host:   host,
		Size:   query.Size,
	}
	r.Processor.Filter = processor.TaskQueueFilter{
		Type:   "img.gen",
		Size:   query.Size,
		Target: query.Prompt,
		API:    query.API,
	}
}

func (web) requestBuilder(c *gin.Context, r *request.Request) {
	var query apiQuery
	var target string
	_ = c.ShouldBindQuery(&query)

	switch {
	case query.Img != "":
		target = query.Img
	case query.URL != "":
		target = query.URL
	}
	urlParse, _ := url.Parse(target)

	referer := c.Request.Referer()
	host := getScheme(c) + c.Request.Host
	ip := c.ClientIP()
	api := urlParse.Host
	device := util.GetDeviceType(c.GetHeader("User-Agent"))
	region := util.GetRegion(ip)

	r.Security.General = security.General{
		Info:    query.More,
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
		UUID:     uuid.New().String(),
		Time:     time.Now(),
		IP:       ip,
		Type:     util.TypeMap["web.img"],
		Target:   target,
		Region:   region,
		Referer:  referer,
		Device:   device,
		API:      api,
		MoreInfo: query.More,
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
