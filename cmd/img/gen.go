package img

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
	"urlAPI/internal/data"
	"urlAPI/internal/file"
	"urlAPI/internal/plugin"
	"urlAPI/internal/security"
)

func GenRequest(IP, Model, API, Target, Size, From, Device string, Referer *url.URL) (ImgResponse, error) {
	var expired = data.Expired
	Domain := Referer.Hostname()
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"img"})))
	if err != nil {
		return ImgResponse{}, err
	}
	if API == "" {
		API = config[0][1]
	}
	if len(config[0]) > 2 {
		expired, err = strconv.Atoi(config[0][2])
		if err != nil {
			return ImgResponse{}, err
		}
	}
	info, err := security.NewRequest(security.SecurityConfig(
		security.WithType("img.gen"),
		security.WithDomain(Domain),
		security.WithAPI(API),
		security.WithIP(IP),
		security.WithTarget(Target)))
	if err != nil {
		return ImgResponse{}, err
	}
	if Model == "" || Size == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{API})))
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
	last, err := data.FetchTask(data.DataConfig(data.WithTaskTarget(Target)))
	if err == nil {
		for _, task := range last {
			if task.Status == "success" && task.Size == Size && task.API == API {
				if time.Now().Sub(task.Time).Minutes() < float64(expired) {
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
						data.WithTaskStatus("outdated")))
					if err != nil {
						log.Println(err)
						return ImgResponse{}, err
					}
				}
			}
		}
	}
	region, err := plugin.GetRegion(plugin.PluginConfig(plugin.WithIP(IP)))
	if err != nil {
		log.Println("Region fetch failed")
	}
	var id string
	if info == "task.except" {
		id = "-1"
	} else {
		id, err = data.NewTask(data.DataConfig(
			data.WithType("图片生成"),
			data.WithAPI(API),
			data.WithTaskIP(IP),
			data.WithTaskTarget(Target),
			data.WithTaskSize(Size),
			data.WithTaskRegion(region.Region),
			data.WithTaskModel(Model),
			data.WithTaskReferer(Referer.String()),
			data.WithTaskDevice(Device),
		))
	}
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
			data.WithTaskStatus("failed"),
			data.WithTaskReturn(err.Error()),
		))
		if editErr != nil {
			err = editErr
		}
		return ImgResponse{}, err
	}
	url := response.URL
	url = strings.ReplaceAll(url, "\\u0026", "&")
	url = strings.ReplaceAll(url, "\\u003c", "<")
	url = strings.ReplaceAll(url, "\\u003e", ">")
	err = file.Add(file.FileConfig(
		file.WithType("img.download"),
		file.WithUUID(id),
		file.WithURL(url),
	))
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
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithTaskStatus("success"),
		data.WithTaskReturn(string(jsonReturn)),
	))
	if err != nil {
		return ImgResponse{}, err
	} else {
		return ret, nil
	}
}
