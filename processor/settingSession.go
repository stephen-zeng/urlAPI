package processor

import (
	"encoding/json"
	"github.com/pkg/errors"
	"urlAPI/database"
)

var PartMap = map[string][]string{
	"openai":       []string{"openai"},
	"deepseek":     []string{"deepseek"},
	"alibaba":      []string{"alibaba"},
	"otherapi":     []string{"otherapi"},
	"security":     []string{"dash", "dashallowedip", "allowedref"},
	"txt":          []string{"txt", "txtgenenabled", "txtacceptprompt"},
	"img":          []string{"img", "imgacceptprompt"},
	"web":          []string{"web", "webimgallowed"},
	"rand":         []string{"rand"},
	"contxt":       []string{"context", "prompt"},
	"taskBehavior": []string{"taskexceptdomain", "taskexceptinfo"},
}

func fetchSetting(info *Session, data *database.Session) error {
	fetchList := PartMap[info.SettingPart]
	info.SettingName = fetchList
	for _, name := range fetchList {
		settingGetter := database.Setting{Name: name}
		settingDBList, err := settingGetter.Read()
		if err != nil {
			return errors.WithStack(err)
		}
		var setting []string
		if err = json.Unmarshal([]byte(settingDBList.SettingList[0].Value), &setting); err != nil {
			return errors.WithStack(err)
		}
		info.SettingData = append(info.SettingData, setting)
	}
	return nil
}

func editSetting(info *Session, data *database.Session) error {
	editList := PartMap[info.SettingPart]
	for index, name := range editList {
		jsonList, err := json.Marshal(info.SettingEdit[index])
		if err != nil {
			return errors.WithStack(err)
		}
		settingEditor := database.Setting{
			Name:  name,
			Value: string(jsonList),
		}
		if err = settingEditor.Update(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
