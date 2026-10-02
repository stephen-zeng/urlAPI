package handles

import (
	"time"
	"urlAPI/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"urlAPI/internal/op"
	"urlAPI/internal/server/middleware"
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

func WebHandler(c *gin.Context) {
	target := c.Query("img")
	referer := c.Request.Referer()
	ip := c.ClientIP()
	task := model.Task{
		UUID:     uuid.New().String(),
		Time:     time.Now(),
		IP:       ip,
		Type:     util.TypeMap["web.img"],
		Target:   target,
		Region:   util.GetRegion(ip),
		Referer:  referer,
		Device:   util.GetDeviceType(c.GetHeader("User-Agent")),
		API:      util.WebTargetHost(target),
		MoreInfo: c.Query("more"),
	}
	_, result, err := op.GenerateWebImage(task, middleware.GetSkipDB(c))
	util.ErrorPrinter(err)
	if err != nil {
		errorReturner(c, err)
		return
	}
	returner(c, result)
}
