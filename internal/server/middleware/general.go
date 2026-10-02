package middleware

import (
	"fmt"
	"sync"
	"time"
	"urlAPI/internal/database"
	"urlAPI/util"
)

type FrequencyFilter struct {
	Type string `json:"type"`
	IP   string `json:"ip"`
}

type FrequencyData struct {
	Counter int       `json:"counter"`
	Time    time.Time `json:"time"`
}

type SafeIPFrequency struct {
	mu          sync.Mutex
	IPFrequency map[FrequencyFilter]FrequencyData
}

var IPFrequency = SafeIPFrequency{
	IPFrequency: make(map[FrequencyFilter]FrequencyData),
}

const (
	frequencyWindow = 250 * time.Millisecond
	// frequencyPruneThreshold bounds the tracker: once it holds this many
	// entries, entries whose window has elapsed are discarded.
	frequencyPruneThreshold = 10000
)

// pruneLocked drops entries whose rate-limit window has expired. Expired
// entries would be reset on their next access anyway, so this does not change
// limiting behaviour. The caller must hold mu.
func (f *SafeIPFrequency) pruneLocked(now time.Time) {
	for filter, data := range f.IPFrequency {
		if now.Sub(data.Time) > frequencyWindow {
			delete(f.IPFrequency, filter)
		}
	}
}

func checkGeneralSecurity(general *General) {
	checkFrequency(general)
	checkException(general)
	checkInfo(general)
}

func checkFrequency(general *General) {
	// 上锁，解锁
	IPFrequency.mu.Lock()
	defer IPFrequency.mu.Unlock()

	filter := FrequencyFilter{
		Type: general.Type,
		IP:   general.IP,
	}
	value, exists := IPFrequency.IPFrequency[filter]
	if !exists {
		if len(IPFrequency.IPFrequency) >= frequencyPruneThreshold {
			IPFrequency.pruneLocked(general.Time)
		}
		value = FrequencyData{}
		value.Counter = 1
		value.Time = time.Now()
		IPFrequency.IPFrequency[filter] = value
		return
	}
	switch {
	case general.Time.Sub(value.Time).Seconds() <= 0.25 && value.Counter >= 10:
		general.Unsafe = true
		general.Info = fmt.Sprintf("%s accessed too frequently", general.IP)
		return
	case general.Time.Sub(value.Time).Seconds() > 0.25:
		value.Counter = 1
		value.Time = time.Now()
	case value.Counter < 10:
		value.Counter++
	}
	IPFrequency.IPFrequency[filter] = value
	return
}

func checkInfo(general *General) {
	if general.Target == "" {
		general.Info = "Empty Target"
		general.Unsafe = true
	}
	allowedref := database.SettingsStore.Get().Security.AllowedReferers
	domain := util.GetDomain(general.Referer)
	if !util.WildcardChecker(&allowedref, &domain) || general.Referer == "" {
		general.Info = fmt.Sprintf("Referer %s not allowed", general.Referer)
		general.Unsafe = true
	}
	return
}

func checkException(general *General) {
	settings := database.SettingsStore.Get()
	taskexceptdomain := settings.Task.ExceptDomains
	taskexceptinfo := settings.Task.ExceptInfos
	domain := util.GetDomain(general.Referer)
	auxInfo := general.Info
	if util.WildcardChecker(&taskexceptdomain, &domain) || util.WildcardChecker(&taskexceptinfo, &auxInfo) {
		general.SkipDB = true
		return
	}
	return
}
