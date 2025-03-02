package txt

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"time"
	"urlAPI/internal/data"
	"urlAPI/internal/plugin"
	"urlAPI/internal/security"
)

func SumRequest(IP, From, Model, API, Target, Device string, Referer *url.URL) (TxtResponse, error) {
	expired := 60
	hash := md5.Sum([]byte(Target))
	md5 := hex.EncodeToString(hash[:])
	Domain := Referer.Hostname()
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"txt"})))
	if err != nil {
		return TxtResponse{}, err
	}
	if API == "" {
		API = config[0][2]
	}
	if len(config[0]) > 3 {
		var err error
		expired, err = strconv.Atoi(config[0][3])
		if err != nil {
			return TxtResponse{}, err
		}
	}
	err = security.NewRequest(security.SecurityConfig(
		security.WithType("txt.sum"),
		security.WithAPI(API),
		security.WithTarget(md5),
		security.WithIP(IP),
		security.WithDomain(Domain),
	))
	if Model == "" {
		config, err = data.FetchSetting(data.DataConfig(data.WithSettingName([]string{API})))
		if err != nil {
			return TxtResponse{}, err
		}
		Model = config[0][2]
	}
	last, err := data.FetchTask(data.DataConfig(data.WithTaskTarget(md5)))
	if err == nil {
		for _, task := range last {
			if task.Status == "success" && task.API == API {
				if time.Now().Sub(task.Time).Minutes() < float64(expired) {
					log.Println("Found old task")
					var ret TxtResponse
					err = json.Unmarshal([]byte(task.Return), &ret)
					if err != nil {
						return TxtResponse{}, err
					} else {
						return ret, nil
					}
				} else {
					log.Println("Found outdated task")
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
		data.WithType("文字总结"),
		data.WithAPI(API),
		data.WithTaskIP(IP),
		data.WithTaskTarget(md5),
		data.WithTaskRegion(region.Region),
		data.WithTaskModel(Model),
		data.WithTaskReferer(Referer.String()),
		data.WithTaskDevice(Device),
	))
	if err != nil {
		return TxtResponse{}, err
	}
	response, err := plugin.Request(plugin.PluginConfig(
		plugin.WithModel(Model),
		plugin.WithAPI(API),
		plugin.WithSumPrompt(Target),
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
		data.WithTaskStatus("success"),
		data.WithTaskReturn(string(jsonRet)),
	))
	if err != nil {
		return TxtResponse{}, err
	} else {
		return ret, nil
	}
}
