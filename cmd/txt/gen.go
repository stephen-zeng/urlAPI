package txt

import (
	"backend/internal/data"
	"backend/internal/plugin"
	"backend/internal/security"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"
)

var shortcut = map[string]string{
	"laugh":    "讲一个笑话，要求在20个字以内，不要换行",
	"poem":     "做诗歌，两句，不要换行。",
	"sentence": "写一句心灵鸡汤，不要换行",
}

func genRequest(IP, Domain, Model, API, Target string) (TxtResponse, error) {
	var target string
	var expired = 60
	if Target == "laugh" || Target == "sentence" || Target == "poem" {
		target = shortcut[Target]
	} else if Target == "" {
		return TxtResponse{}, errors.New("prompt required")
	} else {
		target = Target
		Target = "other"
	}
	config, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"txt"})))
	if err != nil {
		return TxtResponse{}, err
	}
	if API == "" {
		API = config[0][1]
	}
	if len(config[0]) > 3 {
		var err error
		expired, err = strconv.Atoi(config[0][3])
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
		config, err := data.FetchSetting(data.DataConfig(data.WithName([]string{API})))
		if err != nil {
			return TxtResponse{}, nil
		}
		Model = config[0][1]
	}
	last, err := data.FetchTask(data.DataConfig(data.WithTarget(target)))
	if err == nil {
		for _, task := range last {
			if time.Now().Sub(task.Time).Minutes() < float64(expired) && task.Status == "success" && task.API == API {
				log.Println("Found old task")
				var ret TxtResponse
				err := json.Unmarshal([]byte(task.Return), &ret)
				if err != nil {
					return TxtResponse{}, err
				} else {
					return ret, nil
				}
			}
		}
	}
	id, err := data.NewTask(data.DataConfig(
		data.WithIP(IP),
		data.WithTarget(target),
		data.WithType("文字生成"),
	))
	if err != nil {
		return TxtResponse{}, err
	}
	response, err := plugin.Request(plugin.PluginConfig(
		plugin.WithModel(Model),
		plugin.WithAPI(API),
		plugin.WithGenPrompt(target)))
	region, err := plugin.GetRegion(plugin.PluginConfig(plugin.WithIP(IP)))
	if err != nil {
		log.Println("Region fetch failed")
	}
	if err != nil {
		editErr := data.EditTask(data.DataConfig(
			data.WithUUID(id),
			data.WithStatus("failed"),
			data.WithAPI(API),
			data.WithTarget(target),
			data.WithRegion(region.Region)))
		if editErr != nil {
			err = editErr
		}
		return TxtResponse{}, err
	}
	ret := TxtResponse{
		Response: response.Response,
		Prompt:   response.InitPrompt,
		Context:  response.Context,
	}
	jsonRet, err := json.Marshal(ret)
	if err != nil {
		return TxtResponse{}, err
	}
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithReturn(string(jsonRet)),
		data.WithStatus("success"),
		data.WithAPI(API),
		data.WithRegion(region.Region)))
	if err != nil {
		return TxtResponse{}, err
	} else {
		return ret, nil
	}
}
