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

func GenRequest(IP, Domain, Model, API, Target, Size, From string) (ImgResponse, error) {
	if API == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"img"})))
		if err != nil {
			return ImgResponse{}, err
		}
		API = config[0][1]
	}
	err := security.NewRequest(security.SecurityConfig(
		security.WithType("img"),
		security.WithDomain(Domain),
		security.WithAPI(API),
		security.WithIP(IP)))
	if err != nil {
		return ImgResponse{}, err
	}
	if Model == "" || Size == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithName([]string{API})))
		if err != nil {
			return ImgResponse{}, err
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
					var ret ImgResponse
					err := json.Unmarshal([]byte(task.Return), &ret)
					if err != nil {
						return ImgResponse{}, err
					} else {
						return ret, nil
					}
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
						return ImgResponse{}, err
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
		return ImgResponse{}, err
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
		return ImgResponse{}, err
	}
	if err != nil {
		log.Println(err)
		return ImgResponse{}, err
	}
	url := response.URL
	url = strings.ReplaceAll(url, "\\u0026", "&")
	url = strings.ReplaceAll(url, "\\u003c", "<")
	url = strings.ReplaceAll(url, "\\u003e", ">")
	err = file.Add(file.FileConfig(
		file.WithType("img"),
		file.WithUUID(id),
		file.WithURL(url)))
	if err != nil {
		return ImgResponse{}, err
	}
	url = From + "/download?img=" + id
	ret := ImgResponse{
		URL:          url,
		InitPrompt:   response.InitPrompt,
		ActualPrompt: response.ActualPrompt,
	}
	jsonReturn, err := json.Marshal(ret)
	if err != nil {
		return ImgResponse{}, err
	}
	region, err := plugin.GetRegion(plugin.PluginConfig(plugin.WithIP(IP)))
	if err != nil {
		log.Println("Region fetch failed")
	}
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithReturn(string(jsonReturn)),
		data.WithStatus("success"),
		data.WithSize(Size),
		data.WithAPI(API),
		data.WithRegion(region.Region)))
	if err != nil {
		return ImgResponse{}, err
	} else {
		return ret, nil
	}
}
