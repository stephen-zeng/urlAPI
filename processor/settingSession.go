package processor

import (
	"encoding/json"
	"urlAPI/database"
)

var PartMap = map[string][]string{
	"openai":       []string{"openai"},
	"deepseek":     []string{"deepseek"},
	"alibaba":      []string{"alibaba"},
	"otherapi":     []string{"otherapi"},
	"security":     []string{"dash", "dashallowedip", "allowedref"},
	"txt":          []string{"txt", "txtgenenabled"},
	"img":          []string{"img"},
	"web":          []string{"web", "webimgallowed"},
	"rand":         []string{"rand"},
	"contxt":       []string{"context", "prompt"},
	"taskBehavior": []string{"taskexceptdomain"},
}

func fetchSetting(info *Session, data *database.Session) error {
	fetchList := PartMap[info.SettingPart]
	info.SettingName = fetchList
	for _, name := range fetchList {
		settingGetter := database.Setting{Name: name}
		settingDBList, err := settingGetter.Read()
		if err != nil {
			return err
		}
		var setting []string
		err = json.Unmarshal([]byte(settingDBList.SettingList[0].Value), &setting)
		info.SettingData = append(info.SettingData, setting)
	}
	return nil
}

func editSetting(info *Session, data *database.Session) error {
	editList := PartMap[info.SettingPart]
	for index, name := range editList {
		jsonList, err := json.Marshal(info.SettingEdit[index])
		if err != nil {
			return err
		}
		settingEditor := database.Setting{
			Name:  name,
			Value: string(jsonList),
		}
		err = settingEditor.Update()
		if err != nil {
			return err
		}
	}
	return nil
}
