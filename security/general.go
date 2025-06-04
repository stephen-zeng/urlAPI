package security

import (
	"fmt"
	"sync"
	"time"
	"urlAPI/database"
	"urlAPI/util"
)

type SafeIPFrequency struct {
	mu          sync.Mutex
	IPFrequency map[FrequencyFilter]FrequencyData
}

var IPFrequency = SafeIPFrequency{
	IPFrequency: make(map[FrequencyFilter]FrequencyData),
}

func (info *General) GeneralChecker() {
	info.FrequencyChecker()
	info.ExceptionChecker()
	info.InfoChecker()
}

func (info *General) FrequencyChecker() {
	// 上锁，解锁
	IPFrequency.mu.Lock()
	defer IPFrequency.mu.Unlock()

	filter := FrequencyFilter{
		Type: info.Type,
		IP:   info.IP,
	}
	value, exists := IPFrequency.IPFrequency[filter]
	if !exists {
		value = FrequencyData{}
		value.Counter = 1
		value.Time = time.Now()
		IPFrequency.IPFrequency[filter] = value
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
	IPFrequency.IPFrequency[filter] = value
	return
}

func (info *General) InfoChecker() {
	if info.Target == "" {
		info.Info = "Empty Target"
		info.Unsafe = true
	}
	allowedref := database.SettingMap["allowedref"]
	domain := util.GetDomain(info.Referer)
	if !util.RegexChecker(&allowedref, &domain) || info.Referer == "" {
		info.Info = fmt.Sprintf("Referer %s not allowed", info.Referer)
		info.Unsafe = true
	}
	return
}

func (info *General) ExceptionChecker() {
	taskexceptdomain := database.SettingMap["taskexceptdomain"]
	taskexceptinfo := database.SettingMap["taskexceptinfo"]
	domain := util.GetDomain(info.Referer)
	auxInfo := info.Info
	if util.RegexChecker(&taskexceptdomain, &domain) || util.RegexChecker(&taskexceptinfo, &auxInfo) {
		info.SkipDB = true
		return
	}
	return
}
