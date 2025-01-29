package txt

import (
	"backend/internal/data"
	"backend/internal/plugin"
	"backend/internal/security"
	"errors"
	"log"
	"time"
)

var shortcut = map[string]string{
	"laugh":    "讲一个笑话，要求在20个字以内，不要换行",
	"poem":     "做诗歌，两句，不要换行。",
	"sentence": "写一句心灵鸡汤，不要换行",
}

func genRequest(IP, Domain, Model, API, Target string) (string, error) {
	var target string
	if Target == "laugh" || Target == "sentence" || Target == "poem" {
		target = shortcut[Target]
	} else if Target == "" {
		return "", errors.New("prompt required")
	} else {
		target = Target
		Target = "other"
	}
	if API == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithName("txt")))
		if err != nil {
			return "", err
		}
		API = config[0][0]
	}
	err := security.NewRequest(security.SecurityConfig(
		security.WithType("gen"),
		security.WithAPI(API),
		security.WithTarget(Target)))
	if err != nil {
		return "", err
	}
	if Model == "" {
		config, err := data.FetchSetting(data.DataConfig(data.WithName(API)))
		if err != nil {
			return "", nil
		}
		Model = config[0][1]
	}
	last, err := data.FetchTask(data.DataConfig(data.WithTarget(target)))
	if err == nil {
		for _, task := range last {
			if time.Now().Sub(task.Time).Minutes() < 10 && task.Status == "success" {
				log.Println("Found old task")
				return task.Return, nil
			}
		}
	}
	id, err := data.NewTask(data.DataConfig(
		data.WithIP(IP),
		data.WithTarget(target),
		data.WithType("文字生成"),
	))
	if err != nil {
		return "", err
	}
	response, err := plugin.Request(plugin.PluginConfig(
		plugin.WithModel(Model),
		plugin.WithAPI(API),
		plugin.WithGenPrompt(target)))
	if err != nil {
		editErr := data.EditTask(data.DataConfig(
			data.WithUUID(id),
			data.WithStatus("failed")))
		if editErr != nil {
			err = editErr
		}
		return "", err
	}
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithReturn(response),
		data.WithStatus("success")))
	if err != nil {
		return "", err
	} else {
		return response, nil
	}
}
