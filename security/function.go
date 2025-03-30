package security

import (
	"fmt"
	"github.com/pkg/errors"
	"urlAPI/database"
	"urlAPI/util"
)

func (info *TxtGen) FunctionChecker(general *General) error {
	txtgenenabled := database.SettingMap["txtgenenabled"]
	var prompt string
	if _, ok := database.PromptMap[info.Target]; ok {
		prompt = info.Target
	} else {
		prompt = "other"
	}
	switch {
	case database.SettingMap["txt"][0] != "true":
		general.Info = "Txt is not enabled"
		break
	case !util.ListChecker(&txtgenenabled, &(prompt)):
		general.Info = fmt.Sprintf("Target %s is not enabled", info.Target)
		break
	default:
		return nil
	}
	general.Unsafe = true
	return errors.WithStack(errors.New(general.Info))
}

func (info *TxtSum) FunctionChecker(general *General) error {
	switch {
	case database.SettingMap["txt"][0] != "true":
		general.Info = "Txt is not enabled"
		break
	default:
		return nil
	}
	general.Unsafe = true
	return errors.WithStack(errors.New(general.Info))
}

func (info *ImgGen) FunctionChecker(general *General) error {
	switch {
	case database.SettingMap["img"][0] != "true":
		general.Info = "Img is not enabled"
		break
	default:
		return nil
	}
	general.Unsafe = true
	return errors.WithStack(errors.New(general.Info))
}

func (info *Rand) FunctionChecker(general *General) error {
	switch {
	case database.SettingMap["rand"][0] != "true":
		general.Info = "Random is not enabled"
		break
	default:
		return nil
	}
	general.Unsafe = true
	return errors.WithStack(errors.New(general.Info))
}

func (info *WebImg) FunctionChecker(general *General) error {
	webimgallowed := database.SettingMap["webimgallowed"]

	switch {
	case database.SettingMap["web"][1] != "true":
		general.Info = "WebImg is not enabled"
		break
	case !util.ListChecker(&webimgallowed, &(info.API)):
		general.Info = fmt.Sprintf("API %s is not enabled", info.API)
		break
	case info.API == "www.ithome.com" && database.SettingMap["txt"][0] != "true":
		general.Info = "For ITHome, TxtSum is not enabled"
		break
	default:
		return nil
	}
	general.Unsafe = true
	return errors.WithStack(errors.New(general.Info))
}
