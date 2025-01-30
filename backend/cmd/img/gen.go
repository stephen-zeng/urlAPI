package img

import (
	"backend/internal/data"
	"backend/internal/file"
	"backend/internal/plugin"
	"backend/internal/security"
	"encoding/json"
	"log"
	"strings"
	"time"
)

func genRequest(IP, Domain, Model, API, Target, Size, From string) (string, error) {
	if API == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithName("img")))
		if err != nil {
			return "", err
		}
		API = config[0][1]
	}
	err := security.NewRequest(security.SecurityConfig(
		security.WithType("img"),
		security.WithDomain(Domain),
		security.WithAPI(API),
		security.WithIP(IP)))
	if err != nil {
		return "", err
	}
	if Model == "" || Size == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithName(API)))
		if err != nil {
			return "", err
		}
		if Model == "" {
			Model = config[0][3]
		}
		if Size == "" {
			Size = config[0][4]
		}
	}
	last, err := data.FetchTask(data.DataConfig(data.WithTarget(Target)))
	if err == nil {
		for _, task := range last {
			if task.Status == "success" && task.Size == Size && task.API == API {
				if time.Now().Sub(task.Time).Minutes() < 10 {
					log.Println("Found old task")
					return task.Return, nil
				} else {
					log.Println("Found outdated task")
					err = file.Del(file.FileConfig(
						file.WithType("img"),
						file.WithUUID(task.UUID)))
					if err != nil {
						log.Println(err)
					}
					err = data.EditTask(data.DataConfig(
						data.WithUUID(task.UUID),
						data.WithStatus("outdated")))
					if err != nil {
						log.Println(err)
						return "", err
					}
				}
			}
		}
	}
	id, err := data.NewTask(data.DataConfig(
		data.WithIP(IP),
		data.WithTarget(Target),
		data.WithType("图片生成")))
	if err != nil {
		return "", err
	}
	response, err := plugin.Request(plugin.PluginConfig(
		plugin.WithModel(Model),
		plugin.WithAPI(API),
		plugin.WithImgPrompt(Target),
		plugin.WithSize(Size)))
	if err != nil {
		editErr := data.EditTask(data.DataConfig(
			data.WithUUID(id),
			data.WithStatus("failed")))
		if editErr != nil {
			err = editErr
		}
		return "", err
	}
	mapRet := make(map[string]interface{})
	err = json.Unmarshal([]byte(response), &mapRet)
	if err != nil {
		log.Println(err)
		return "", err
	}
	url := mapRet["url"].(string)
	url = strings.ReplaceAll(url, "\\u0026", "&")
	url = strings.ReplaceAll(url, "\\u003c", "<")
	url = strings.ReplaceAll(url, "\\u003e", ">")
	err = file.Add(file.FileConfig(
		file.WithType("img"),
		file.WithUUID(id),
		file.WithURL(url)))
	if err != nil {
		return "", err
	}
	url = From + "/download?img=" + id
	mapRet["url"] = url
	jsonRet, err := json.Marshal(mapRet)
	if err != nil {
		return "", err
	}
	ret := string(jsonRet)
	ret = strings.ReplaceAll(ret, "\\u0026", "&")
	ret = strings.ReplaceAll(ret, "\\u003c", "<")
	ret = strings.ReplaceAll(ret, "\\u003e", ">")
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithReturn(ret),
		data.WithStatus("success"),
		data.WithSize(Size),
		data.WithAPI(API)))
	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}
