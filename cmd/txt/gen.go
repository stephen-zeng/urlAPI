package txt

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"
	"urlAPI/cmd/img"
	"urlAPI/internal/data"
	"urlAPI/internal/file"
	"urlAPI/internal/plugin"
	"urlAPI/internal/security"
)

func GenRequest(IP, From, Model, API, Target string, Referer *url.URL) (TxtResponse, error) {
	var target string
	var expired = 60
	Domain := Referer.Hostname()
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"prompt", "txt"})))
	switch Target {
	case "laugh":
		target = config[0][0]
	case "poem":
		target = config[0][1]
	case "sentence":
		target = config[0][2]
	case "":
		return TxtResponse{}, errors.New("prompt required")
	default:
		target = Target
		Target = "other"
	}
	if API == "" {
		API = config[1][1]
	}
	if len(config[1]) > 3 {
		var err error
		expired, err = strconv.Atoi(config[1][3])
		if err != nil {
			return TxtResponse{}, err
		}
	}
	err = security.NewRequest(security.SecurityConfig(
		security.WithType("txt.gen"),
		security.WithAPI(API),
		security.WithTarget(Target),
		security.WithDomain(Domain),
		security.WithIP(IP),
		security.WithTarget(Target)))
	if err != nil {
		return TxtResponse{}, err
	}
	if Model == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{API})))
		if err != nil {
			return TxtResponse{}, nil
		}
		Model = config[0][1]
	}
	last, err := data.FetchTask(data.DataConfig(data.WithTaskTarget(target)))
	if err == nil {
		for _, task := range last {
			if task.Status == "success" && task.API == API {
				if time.Now().Sub(task.Time).Minutes() < float64(expired) {
					log.Println("Found old task")
					var ret TxtResponse
					err := json.Unmarshal([]byte(task.Return), &ret)
					if err != nil {
						return TxtResponse{}, err
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
						return TxtResponse{}, err
					}
				}
			}
		}
	}
	region, err := plugin.GetRegion(plugin.PluginConfig(plugin.WithIP(IP)))
	if err != nil {
		log.Println("Region fetch failed")
	}
	id, err := data.NewTask(data.DataConfig(
		data.WithType("文字生成"),
		data.WithAPI(API),
		data.WithTaskIP(IP),
		data.WithTaskTarget(target),
		data.WithTaskRegion(region.Region),
		data.WithTaskModel(Model),
		data.WithTaskReferer(Referer.String()),
	))
	if err != nil {
		return TxtResponse{}, err
	}
	response, err := plugin.Request(plugin.PluginConfig(
		plugin.WithModel(Model),
		plugin.WithAPI(API),
		plugin.WithGenPrompt(target),
	))
	if err != nil {
		editErr := data.EditTask(data.DataConfig(
			data.WithUUID(id),
			data.WithTaskStatus("failed"),
			data.WithTaskReturn(err.Error()),
		))
		if editErr != nil {
			err = editErr
		}
		return TxtResponse{}, err
	}
	imgResponse, err := img.TxtDrawRequest(id, response.Response, From)
	if err != nil {
		log.Println(err)
	}
	ret := TxtResponse{
		Response: response.Response,
		Prompt:   response.InitPrompt,
		Context:  response.Context,
		URL:      imgResponse.URL,
	}
	jsonRet, err := json.Marshal(ret)
	if err != nil {
		return TxtResponse{}, err
	}
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithTaskStatus("success"),
		data.WithTaskReturn(string(jsonRet)),
	))
	if err != nil {
		return TxtResponse{}, err
	} else {
		return ret, nil
	}
}
