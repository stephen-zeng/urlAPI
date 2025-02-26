package web

import (
	"encoding/json"
	"errors"
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

var webMap = map[string]string{
	"github.com":       "github",
	"gitee.com":        "gitee",
	"www.bilibili.com": "bilibili",
	"www.youtube.com":  "youtube",
	"arxiv.org":        "arxiv",
	"www.ithome.com":   "ithome",
}

func getBiliABV(URL string) string {
	for i := 31; i < len(URL); i++ {
		if URL[i] == '/' {
			return URL[31:i]
		}
	}
	return ""
}
func getYtbID(URL string) string {
	for i := 32; i < len(URL); i++ {
		if URL[i] == '&' {
			return URL[32:i]
		}
	}
	return URL[32:]
}

func ImgRequest(IP, From, API, Target string, Referer *url.URL) (WebResponse, error) {
	var expired = data.Expired
	Domain := Referer.Hostname()
	API = webMap[API] // github.com --> github
	err := security.NewRequest(security.SecurityConfig(
		security.WithType("web.img"),
		security.WithAPI(API),
		security.WithDomain(Domain),
		security.WithIP(IP)))
	if err != nil {
		return WebResponse{}, err
	}
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"web"})))
	if err != nil {
		return WebResponse{}, err
	}
	if len(config[0]) > 3 {
		var err error
		expired, err = strconv.Atoi(config[0][3])
		if err != nil {
			return WebResponse{}, err
		}
	}
	last, err := data.FetchTask(data.DataConfig(data.WithTaskTarget(Target)))
	if err == nil {
		for _, task := range last {
			if task.Status == "success" && task.API == API {
				if time.Now().Sub(task.Time).Minutes() < float64(expired) {
					log.Println("Found old task")
					var ret WebResponse
					err := json.Unmarshal([]byte(task.Return), &ret)
					if err != nil {
						return WebResponse{}, err
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
						return WebResponse{}, err
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
		data.WithType("网站缩略图"),
		data.WithAPI(API),
		data.WithTaskIP(IP),
		data.WithTaskTarget(Target),
		data.WithTaskRegion(region.Region),
		data.WithTaskReferer(Referer.String()),
	))
	var ret WebResponse
	var token string
	switch API {
	case "github":
		if len(config[0]) > 5 {
			token = config[0][5]
		}
		ret, err = repo(strings.ReplaceAll(Target, "https://github.com", "https://api.github.com/repos"), From, id, token)
	case "gitee":
		ret, err = repo(strings.ReplaceAll(Target, "https://gitee.com", "https://gitee.com/api/v5/repos"), From, id, token)
	case "bilibili":
		ret, err = bili(getBiliABV(Target), From, id)
	case "youtube":
		if len(config[0]) > 6 {
			ret, err = ytb(getYtbID(Target), From, id, config[0][6])
		} else {
			err = errors.New("No Youtube API Token")
		}
	case "arxiv":
		ret, err = arxiv(Target, From, id)
	case "ithome":
		ret, err = ithome(Target, From, id, IP, Referer)
	default:
		err = errors.New("Unsupported websites")
	}
	if err != nil {
		editErr := data.EditTask(data.DataConfig(
			data.WithUUID(id),
			data.WithTaskStatus("failed"),
			data.WithTaskReturn(err.Error()),
		))
		if editErr != nil {
			err = editErr
		}
		return WebResponse{}, err
	}
	ret.Target = Target
	jsonRet, err := json.Marshal(ret)
	if err != nil {
		return WebResponse{}, err
	}
	err = data.EditTask(data.DataConfig(
		data.WithUUID(id),
		data.WithTaskStatus("success"),
		data.WithTaskReturn(string(jsonRet)),
	))
	if err != nil {
		return WebResponse{}, err
	} else {
		return ret, nil
	}
}
