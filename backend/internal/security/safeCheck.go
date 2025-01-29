package security

import (
	"backend/internal/data"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

func dashCheck(IP, Pwd string) error {
	list, err := data.FetchSetting(data.DataConfig(data.WithName("dash")))
	if err != nil {
		log.Println(err)
		return err
	}
	for index, item := range list[0] {
		if index == 0 {
			if item != Pwd {
				return errors.New("Dashboard Passwords don't match.")
			} else {
				rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(item), `\*`, ".*") + "$"
				match, err := regexp.MatchString(rgx, IP)
				if err != nil {
					continue
				}
				if match {
					return nil
				}
			}
		}
	}
	log.Println("The IP " + IP + " is NOT permitted to access the dashboard.")
	return errors.New("Not in the dashboard IP Whitelist")
}

var Frequency map[string]time.Time = make(map[string]time.Time)

func frequencyCheck(IP string) error {
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
			log.Println("The IP " + IP + " accessed too frequent.")
			return errors.New("frequencyCheck failed")
		}
		return nil
	}
}

func sourceCheck(source string) error {
	list, err := data.FetchSetting(data.DataConfig(data.WithName("allow")))
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

func webTargetCheck(target string) error {
	list, err := data.FetchSetting(data.DataConfig(data.WithName("blocklist")))
	if err != nil {
		return err
	}
	for _, item := range list[0] {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(item), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, target)
		if err != nil {
			continue
		}
		if match {
			log.Println("Target " + target + " is in blacklist.")
			return errors.New("targetCheck failed")
		}
	}
	return nil
}

func txtGenTargetCheck(Target string) error {
	txtEnabled, err := data.FetchSetting(data.DataConfig(data.WithName("txtgenenabled")))
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
