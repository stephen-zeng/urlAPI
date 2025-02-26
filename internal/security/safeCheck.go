package security

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
	"urlAPI/internal/data"
)

var typeMap = map[string]string{
	"txt.gen":  "文字生成",
	"txt.sum":  "文字总结",
	"img.gen":  "图片生成",
	"download": "文件下载",
	"rand":     "随机图片",
	"web.img":  "网站缩略图",
}

func frequencyCheck(IP, Type, Target string) error {
	if Type == "rand" {
		return nil
	}
	list, err := data.FetchTask(data.DataConfig(data.WithTaskIP(IP)))
	if err != nil {
		if err.Error() == "Task not found" {
			return nil
		} else {
			return err
		}
	}
	for _, task := range list {
		if time.Now().Sub(task.Time).Seconds() <= 0.25 && task.Type == typeMap[Type] && task.Target == Target {
			log.Println("The IP " + IP + " accessed too frequently.")
			return errors.New("frequencyCheck failed")
		}
	}
	return nil
}

func sourceCheck(source string) error {
	list, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"allowedref"})))
	if err != nil {
		return err
	}
	for _, item := range list[0] {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(item), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, source)
		if err != nil {
			continue
		}
		if match {
			return nil
		}
	}
	log.Println("Source " + source + " is Not in whitelist.")
	return errors.New("sourceCheck failed")
}

func txtCheck(Target, Type string) error {
	if Target == "" {
		log.Println("The target is empty.")
		return errors.New("txtTargetCheck failed")
	}
	txt, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"txt", "txtgenenabled"})))
	if err != nil {
		return err
	}
	if txt[0][0] != "true" {
		log.Println("Txt disabled.")
		return errors.New("Txt disabled")
	}
	if Type == "txt.gen" {
		for _, item := range txt[1] {
			if item == Target {
				return nil
			}
		}
		log.Println("The target " + Target + " is NOT enabled.")
		return errors.New("txtGenCheck failed")
	} else {
		return nil
	}
}

func imgGenCheck() error {
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"img"})))
	if err != nil {
		return err
	}
	if config[0][0] == "true" {
		return nil
	} else {
		log.Println("The imgGen isn't enabled.")
		return errors.New("Img Disabled")
	}
}

func randCheck(API, Info string) error {
	rd, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"rand"})))
	if err != nil {
		return err
	}
	if rd[0][0] != "true" {
		log.Println("Random picture isn't enabled.")
		return errors.New("Random picture disabled")
	}
	_, err = data.FetchRepo(data.DataConfig(
		data.WithBy("api&info"),
		data.WithAPI(API),
		data.WithRepoInfo(Info)))
	if err != nil {
		return err
	} else {
		return nil
	}
}

func webImgCheck(API string) error {
	list, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"web", "webimgallowed"})))
	if err != nil {
		return err
	}
	if list[0][1] != "true" {
		log.Println("WebImg isn't enabled.")
		return errors.New("WebImg Disabled")
	}
	for _, item := range list[1] {
		if item == API {
			return nil
		}
	}
	log.Println("The Website " + API + " is NOT enabled for webimg.")
	return errors.New("webImgCheck failed")
}
