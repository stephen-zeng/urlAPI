package security

import (
	"fmt"
	"time"
	"urlAPI/database"
	"urlAPI/util"
)

var IPFrequency = make(map[FrequencyFilter]FrequencyData)

func (info *General) GeneralChecker() {
	info.FrequencyChecker()
	info.ExceptionChecker()
	info.InfoChecker()
}

func (info *General) FrequencyChecker() {
	filter := FrequencyFilter{
		Type: info.Type,
		IP:   info.IP,
	}
	value, exists := IPFrequency[filter]
	if !exists {
		value = FrequencyData{}
		value.Counter = 1
		value.Time = time.Now()
		IPFrequency[filter] = value
		return
	}
	switch {
	case info.Time.Sub(value.Time).Seconds() <= 0.25 && value.Counter >= 10:
		info.Unsafe = true
		info.Info = fmt.Sprintf("%s accessed too frequently", info.IP)
		return
	case info.Time.Sub(value.Time).Seconds() > 0.25:
		value.Counter = 1
		value.Time = time.Now()
	case value.Counter < 10:
		value.Counter++
	}
	IPFrequency[filter] = value
	return
}

func (info *General) InfoChecker() {
	if info.Target == "" {
		info.Info = "Empty Target"
		info.Unsafe = true
	}
	allowedref := database.SettingMap["allowedref"]
	domain := util.GetDomain(info.Referer)
	if !util.WildcardChecker(&allowedref, &domain) || info.Referer == "" {
		info.Info = fmt.Sprintf("Referer %s not allowed", info.Referer)
		info.Unsafe = true
	}
	return
}

func (info *General) ExceptionChecker() {
	taskexceptdomain := database.SettingMap["taskexceptdomain"]
	domain := util.GetDomain(info.Referer)
	if util.WildcardChecker(&taskexceptdomain, &domain) {
		info.SkipDB = true
		return
	}
	return
}
