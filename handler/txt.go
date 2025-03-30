package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/security"
	"urlAPI/util"
)

func txtHandler(c *gin.Context) {
	var txtRequest request.Request
	txtRequestBuilder(c, &txtRequest)

	if err := txtChecker(&txtRequest); err != nil {
		log.Printf("%s from %s\n", err, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": txtRequest.Security.General.Info,
		})
		return
	}
	if txtOldTask(&txtRequest) {
		returner(c, txtRequest.DB.Task.Return, txtRequest.Processor.TxtGen.Return)
		return
	}
	util.ErrorPrinter(txtRequest.Processor.TxtGen.Process(&txtRequest.DB.Task))
	taskSaver(&txtRequest)
	returner(c, txtRequest.DB.Task.Return, txtRequest.Processor.TxtGen.Return)
	return
}

func txtOldTask(r *request.Request) bool {
	var hasOldTask bool
	expireTime, _ := strconv.Atoi(database.SettingMap["txt"][3])
	taskFinder := database.Task{
		Type:   r.DB.Task.Type,
		Target: r.DB.Task.Target,
		Status: "success",
	}
	taskDBList, err := taskFinder.Read()
	if err != nil {
		log.Printf("Handler txtHandler %s", err.Error())
		return false
	}
	for _, task := range taskDBList.TaskList {
		if time.Now().Sub(task.Time) <= time.Duration(expireTime)*time.Minute {
			hasOldTask = true
			r.Processor.TxtGen.Return = r.Processor.TxtGen.Host + "/download?img=" + task.UUID
			r.DB.Task.Return = task.Return
		} else {
			os.Remove(processor.ImgPath + task.UUID + ".png")
		}
	}
	return hasOldTask
}

func txtChecker(r *request.Request) error {
	var err error
	err = r.Security.General.FrequencyChecker()
	err = r.Security.General.ExceptionChecker()
	err = r.Security.General.InfoChecker()
	err = r.Security.TxtGen.APIChecker(&r.Security.General)
	err = r.Security.TxtGen.FunctionChecker(&r.Security.General)
	return err
}

func txtRequestBuilder(c *gin.Context, r *request.Request) {
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
	r.Security.TxtSum = security.TxtSum{
		Model: model,
		API:   api,
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
}
