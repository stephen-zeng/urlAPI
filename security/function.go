package security

import (
	"fmt"
	"urlAPI/database"
	"urlAPI/util"
)

func (info *TxtGen) FunctionChecker(general *General) {
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
		return
	}
	general.Unsafe = true
}

func (info *TxtSum) FunctionChecker(general *General) {
	switch {
	case database.SettingMap["txt"][0] != "true":
		general.Info = "Txt is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *ImgGen) FunctionChecker(general *General) {
	switch {
	case database.SettingMap["img"][0] != "true":
		general.Info = "Img is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *Rand) FunctionChecker(general *General) {
	switch {
	case database.SettingMap["rand"][0] != "true":
		general.Info = "Random is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *WebImg) FunctionChecker(general *General) {
	webimgallowed := database.SettingMap["webimgallowed"]

	switch {
	case database.SettingMap["web"][0] != "true":
		general.Info = "WebImg is not enabled"
		break
	case !util.ListChecker(&webimgallowed, &(info.API)):
		general.Info = fmt.Sprintf("API %s is not enabled", info.API)
		break
	case info.API == "www.ithome.com" && database.SettingMap["txt"][0] != "true":
		general.Info = "For ITHome, TxtSum is not enabled"
		break
	default:
		return
	}
	general.Unsafe = true
}

func (info *General) FunctionChecker(general *General) {}
