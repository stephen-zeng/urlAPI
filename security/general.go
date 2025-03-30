package security

import (
	"fmt"
	"github.com/pkg/errors"
	"urlAPI/database"
	"urlAPI/util"
)

var IPFrequency = make(map[string]General)

func (info *General) FrequencyChecker() error {
	value, exists := IPFrequency[info.IP]
	if !exists {
		IPFrequency[info.IP] = *info
		return nil
	}
	if info.Time.Sub(value.Time).Seconds() <= 0.25 && value.Type == info.Type {
		info.Unsafe = true
		info.Info = fmt.Sprintf("%s accessed too frequently", info.IP)
		return errors.WithStack(errors.New(info.Info))
	}
	return nil
}

func (info *General) InfoChecker() error {
	if info.Target == "" {
		info.Info = "Empty Target"
		info.Unsafe = true
		return errors.WithStack(errors.New(info.Info))
	}
	allowedref := database.SettingMap["allowedref"]
	if !util.RefererChecker(&allowedref, &(info.Referer)) || info.Referer == "" {
		info.Info = fmt.Sprintf("Referer %s not allowed", info.Referer)
		info.Unsafe = true
		return errors.WithStack(errors.New(info.Info))
	}
	return nil
}

func (info *General) ExceptionChecker() error {
	taskexceptdomain := database.SettingMap["taskexceptdomain"]
	if util.RefererChecker(&taskexceptdomain, &(info.Referer)) {
		info.SkipDB = true
		return nil
	}
	return nil
}
