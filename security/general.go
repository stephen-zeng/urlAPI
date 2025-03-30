package security

import (
	"fmt"
	"urlAPI/database"
	"urlAPI/util"
)

var IPFrequency = make(map[string]General)

func (info *General) FrequencyChecker() {
	value, exists := IPFrequency[info.IP]
	if !exists {
		IPFrequency[info.IP] = *info
		return
	}
	if info.Time.Sub(value.Time).Seconds() <= 0.25 && value.Type == info.Type {
		info.Unsafe = true
		info.Info = fmt.Sprintf("%s accessed too frequently", info.IP)
	}
	return
}

func (info *General) InfoChecker() {
	if info.Target == "" {
		info.Info = "Empty Target"
		info.Unsafe = true
	}
	allowedref := database.SettingMap["allowedref"]
	if !util.RefererChecker(&allowedref, &(info.Referer)) || info.Referer == "" {
		info.Info = fmt.Sprintf("Referer %s not allowed", info.Referer)
		info.Unsafe = true
	}
	return
}

func (info *General) ExceptionChecker() {
	taskexceptdomain := database.SettingMap["taskexceptdomain"]
	if util.RefererChecker(&taskexceptdomain, &(info.Referer)) {
		info.SkipDB = true
		return
	}
	return
}
