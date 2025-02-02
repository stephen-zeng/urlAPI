package security

import (
	"backend/internal/data"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

var Frequency = make(map[string]time.Time)

func frequencyCheck(IP string) error {
	if IP == "" {
		log.Println("IP is empty")
		return errors.New("frequencyCheck failed")
	}
	if _, exist := Frequency[IP]; !exist {
		if len(Frequency) >= 1000000 {
			Frequency = make(map[string]time.Time)
		}
		Frequency[IP] = time.Now()
		return nil
	} else {
		lastTime := Frequency[IP]
		Frequency[IP] = time.Now()
		if time.Now().Sub(lastTime).Seconds() < 0.25 {
			log.Println("The IP " + IP + " accessed too frequently.")
			return errors.New("frequencyCheck failed")
		}
		return nil
	}
}

func sourceCheck(source string) error {
	if source == "" {
		log.Println("The source domain is empty.")
		return errors.New("sourceCheck failed")
	}
	list, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"allowedref"})))
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

func txtGenCheck(Target string) error {
	if Target == "" {
		log.Println("The target domain is empty.")
		return errors.New("txtGenTargetCheck failed")
	}
	txtEnabled, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"txtgenenabled"})))
	if err != nil {
		return err
	}
	for _, item := range txtEnabled[0] {
		if item == Target {
			return nil
		}
	}
	log.Println("The target " + Target + " is NOT enabled.")
	return errors.New("txtGenCheck failed")
}

func imgGenCheck() error {
	config, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"img"})))
	if err != nil {
		return err
	}
	if config[0][0] == "true" {
		return nil
	} else {
		log.Println("The imgGen isn't enabled.")
		return errors.New("imgGenCheck failed")
	}
}
