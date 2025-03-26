package processor

import "urlAPI/database"

func (info *TxtGen) Process(data *database.Task) error {
	info.Return = database.SettingMap["txt"][4]
	if _, ok := database.PromptMap[info.Target]; ok {
		info.Target = database.SettingMap["prompt"][database.PromptMap[info.Target]]
	}
}
