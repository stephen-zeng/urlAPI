package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
	"urlAPI/database"
	"urlAPI/processor"
	"urlAPI/request"
	"urlAPI/util"
)

type JsonStruct struct {
	Prompt         string `json:"prompt"`
	OriginalPrompt string `json:"original_prompt"`
	ActualPrompt   string `json:"actual_prompt"`
	Response       string `json:"response"`
	URL            string `json:"url"`
}

/*
任务执行后的操作
1. 更新任务列表
2. 写入数据库
*/
func afterTask(r *request.Request) {
	// 运行后更新
	processor.TaskCounter[r.Processor.Filter.API]--
	taskTmp := processor.TaskQueue[r.Processor.Filter]
	taskTmp.Running = false
	if r.DB.Task.Status == "success" {
		taskTmp.Time = time.Now()
		taskTmp.Return = r.Processor.Return
		taskTmp.DBReturn = r.DB.Task.Return
		taskTmp.UUID = r.DB.Task.UUID
	}
	processor.TaskQueue[r.Processor.Filter] = taskTmp
	if !r.Security.General.SkipDB {
		util.ErrorPrinter(r.DB.Task.Create())
	}
}

/*
任务执行前的操作
1. 查找已有任务
2. 等待多余的任务执行完毕
3. 将任务添加至任务队列中
*/
func beforeTask(r *request.Request) {
	settingName := task2settingName[r.Processor.Filter.Type]
	expiredPosition := expiredSettingPosition[settingName]
	expiredTime, _ := strconv.Atoi(database.SettingMap[settingName][expiredPosition])
	if task, ok := processor.TaskQueue[r.Processor.Filter]; ok {
		for {
			if !task.Running {
				break
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(1 * time.Millisecond)
		if !task.Running && time.Now().Sub(task.Time) <= time.Duration(expiredTime)*time.Minute {
			r.Processor.Return = task.Return
			r.DB.Task.Return = task.DBReturn
			r.DB.Task.Status = "success"
			return
		} else {
			os.Remove(processor.ImgPath + task.UUID + ".png")
		}
	}
	// 没有已知任务，准备开新任务
	for {
		value, ok := processor.TaskCounter[r.Processor.Filter.API]
		if !ok || value <= 2 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 运行前准备
	taskTmp := processor.TaskQueue[r.Processor.Filter]
	taskTmp.Running = true
	processor.TaskQueue[r.Processor.Filter] = taskTmp
	processor.TaskCounter[r.Processor.Filter.API]++
}

func returner(c *gin.Context, r *request.Request) {
	var jsonStruct JsonStruct
	json.Unmarshal([]byte(r.DB.Task.Return), &jsonStruct)
	if c.Query("format") == "json" {
		c.JSON(http.StatusOK, jsonStruct)
	} else {
		c.Redirect(http.StatusFound, r.Processor.Return)
	}
}

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return `https://`
	}
	if scheme := c.GetHeader("X-Forwarded-Proto"); scheme != "" {
		return scheme + `://`
	}
	return `http://`
}
