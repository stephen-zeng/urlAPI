package txt

import (
	"backend/cmd/img"
	"backend/internal/data"
	"backend/internal/file"
	"backend/internal/plugin"
	"backend/internal/security"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"
)

var shortcut = map[string]string{
	"laugh":    "讲一个笑话，中文，不要换行，需要句中有标点符号",
	"poem":     "做几句诗歌，中文，不要换行，需要句中有标点符号",
	"sentence": "写几句心灵鸡汤，中文，不要换行，需要句中有标点符号",
}

func GenRequest(IP, From, Domain, Model, API, Target, Regen string) (TxtResponse, error) {
	var target string
	var expired = 60
	var fallbackURL = "https://raw.githubusercontent.com/stephen-zeng/img/master/fallback.png"
	if Target == "laugh" || Target == "sentence" || Target == "poem" {
		target = shortcut[Target]
	} else if Target == "" {
		return TxtResponse{}, errors.New("prompt required")
	} else {
		target = Target
		Target = "other"
	}
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"txt"})))
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
	if len(config[0]) > 4 {
		fallbackURL = config[0][4]
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
	if err == nil && Regen != "true" {
		for _, task := range last {
			if task.Status == "success" && task.API == API {
				if time.Now().Sub(task.Time).Minutes() < float64(expired) {
					log.Println("Found old task")
					var ret TxtResponse
					err := json.Unmarshal([]byte(task.Return), &ret)
					if err != nil {
						return TxtResponse{}, err
					} else if Regen != "true" {
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
		return TxtResponse{
			URL: fallbackURL,
		}, err
	}
	imgResponse, err := img.DrawRequest(id, response.Response, From)
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
