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
	// 读
	processor.TaskCounter.Mu.Lock()
	processor.TaskCounter.Counter[r.Processor.Filter.API]--
	processor.TaskCounter.Mu.Unlock()

	//写
	processor.TaskCounter.Mu.RLock()
	taskTmp := processor.TaskQueue.Queue[r.Processor.Filter]
	processor.TaskCounter.Mu.RUnlock()

	taskTmp.Running = false
	if r.DB.Task.Status == "success" {
		taskTmp.DB = r.DB.Task
		taskTmp.Return = r.Processor.Return
	}

	processor.TaskCounter.Mu.Lock()
	processor.TaskQueue.Queue[r.Processor.Filter] = taskTmp
	processor.TaskCounter.Mu.Unlock()

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
	//正式开始
	settingName := task2settingName[r.Processor.Filter.Type]
	expiredPosition := expiredSettingPosition[settingName]
	expiredTime, _ := strconv.Atoi(database.SettingMap[settingName][expiredPosition])

	processor.TaskCounter.Mu.RLock()
	task, ok := processor.TaskQueue.Queue[r.Processor.Filter]
	processor.TaskCounter.Mu.RUnlock()

	if ok {
		for {
			if !task.Running {
				break
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(1 * time.Millisecond)
		if !task.Running && time.Now().Sub(task.DB.Time) <= time.Duration(expiredTime)*time.Minute {
			r.Processor.Return = task.Return
			id := r.DB.Task.UUID
			r.DB.Task = task.DB
			r.DB.Task.UUID = id
			r.DB.Task.Time = time.Now()
			r.DB.Task.Temp = "Yes"
			r.Processor.Return = task.Return
			return
		} else {
			os.Remove(processor.ImgPath + task.DB.UUID + ".png")
		}
	}

	// 没有已知任务，准备开新任务
	r.DB.Task.Temp = "No"
	for {

		processor.TaskCounter.Mu.RLock()
		value, ok := processor.TaskCounter.Counter[r.Processor.Filter.API]
		processor.TaskCounter.Mu.RUnlock()

		if !ok || value <= 2 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 运行前准备
	processor.TaskCounter.Mu.RLock()
	taskTmp := processor.TaskQueue.Queue[r.Processor.Filter]
	processor.TaskCounter.Mu.RUnlock()

	taskTmp.Running = true

	processor.TaskCounter.Mu.Lock()
	processor.TaskQueue.Queue[r.Processor.Filter] = taskTmp
	processor.TaskCounter.Counter[r.Processor.Filter.API]++
	processor.TaskCounter.Mu.Unlock()
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
