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

func imgHandler(c *gin.Context) {
	var imgRequest request.Request
	imgBuilder(c, &imgRequest)
	if err := imgChecker(&imgRequest); err != nil {
		log.Printf("%s from %s\n", err, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{
			"error": imgRequest.Security.General.Info,
		})
		return
	}
	if imgOldTask(&imgRequest) {
		returner(c, imgRequest.DB.Task.Return, imgRequest.Processor.ImgGen.Return)
		return
	}
	util.ErrorPrinter(imgRequest.Processor.ImgGen.Process(&imgRequest.DB.Task))
	taskSaver(&imgRequest)
	returner(c, imgRequest.DB.Task.Return, imgRequest.Processor.ImgGen.Return)
	return
}

func imgOldTask(r *request.Request) bool {
	var hasOldTask bool
	expireTime, _ := strconv.Atoi(database.SettingMap["img"][2])
	taskFinder := database.Task{
		Target: r.DB.Task.Target,
		Type:   r.DB.Task.Type,
		Status: "success",
	}
	taskDBList, err := taskFinder.Read()
	if err != nil {
		log.Printf("Handler imgHandler %s", err.Error())
		return false
	}
	for _, task := range taskDBList.TaskList {
		if time.Now().Sub(task.Time) <= time.Duration(expireTime)*time.Minute {
			hasOldTask = true
			r.Processor.ImgGen.Return = r.Processor.ImgGen.Host + "/download?img=" + task.UUID
			r.DB.Task.Return = task.Return
		} else {
			os.Remove(processor.ImgPath + task.UUID + ".png")
		}
	}
	return hasOldTask
}

func imgChecker(r *request.Request) error {
	var err error
	err = r.Security.General.ExceptionChecker()
	err = r.Security.General.InfoChecker()
	err = r.Security.General.FrequencyChecker()
	err = r.Security.ImgGen.APIChecker(&r.Security.General)
	err = r.Security.ImgGen.FunctionChecker(&r.Security.General)
	return err
}

func imgBuilder(c *gin.Context, r *request.Request) {
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
}
