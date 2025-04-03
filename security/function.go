package security

import (
	"fmt"
	"urlAPI/database"
	"urlAPI/util"
)

func (info *TxtGen) FunctionChecker(general *General) {
	txtgenenabled := database.SettingMap["txtgenenabled"]
	txtacceptprompt := database.SettingMap["txtacceptprompt"]
	var prompt string
	if _, ok := database.PromptMap[general.Target]; ok {
		prompt = general.Target
	} else {
		prompt = "other"
	}
	switch {
	case database.SettingMap["txt"][0] != "true":
		general.Info = "Txt is not enabled"
		break
	case !util.ListChecker(&txtgenenabled, &(prompt)) || general.Target == "" || !util.WildcardChecker(&txtacceptprompt, &(general.Target)):
		general.Info = fmt.Sprintf("Target %s is not enabled for Txt	Gen", general.Target)
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
	imgacceptprompt := database.SettingMap["imgacceptprompt"]
	switch {
	case database.SettingMap["img"][0] != "true":
		general.Info = "Img is not enabled"
		break
	case general.Target == "" || !util.WildcardChecker(&imgacceptprompt, &(general.Target)):
		general.Info = fmt.Sprintf("Prompt %s is not allowed for ImgGen", general.Target)
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
		return
	}
	general.Unsafe = true
}
