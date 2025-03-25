package security

import (
	"fmt"
	"time"
	"urlAPI/database"
	"urlAPI/util"
)

var IPFrequency = map[string]General{}

func (info *General) FrequencyChecker() {
	value, exists := IPFrequency[info.IP]
	if !exists {
		IPFrequency[info.IP] = *info
		return
	}
	if time.Now().Sub(info.Time).Seconds() <= 0.25 && value.Type == info.Type {
		info.Unsafe = true
		info.Info = fmt.Sprintf("%s accessed too frequently", info.IP)
	}
}

func (info *General) RefererChecker() {
	allowedref := database.SettingMap["allowedref"]
	if !util.RefererChecker(&allowedref, &(info.Referer)) || info.Referer == "" {
		info.Info = "Referer not allowed"
		info.Unsafe = true
	}
}

func (info *General) ExceptionChecker() {
	taskexceptdomain := database.SettingMap["taskexceptdomain"]
	if util.RefererChecker(&taskexceptdomain, &(info.Referer)) {
		info.SkipDB = true
	}
}
