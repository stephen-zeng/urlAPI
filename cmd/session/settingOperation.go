package session

import (
	"backend/internal/data"
)

var PartMap = map[string][]string{
	"openai":   []string{"openai"},
	"deepseek": []string{"deepseek"},
	"alibaba":  []string{"alibaba"},
	"otherapi": []string{"otherapi"},
	"security": []string{"dash", "dashallowedip", "allowedref"},
	"txt":      []string{"txt", "txtgenenabled", "txtsumenabled"},
	"img":      []string{"img"},
	"template": []string{"template", "webimgallowed", "websumblocked"},
	"rand":     []string{"rand"},
}

func fetchSetting(Part string) ([]string, [][]string, error) {
	response, err := data.FetchSetting(data.DataConfig(data.WithSettingName(PartMap[Part])))
	if err != nil {
		return nil, nil, err
	} else {
		return PartMap[Part], response, nil
	}
}

func editSetting(Part string, Edit [][]string) ([]string, error) {
	err := data.EditSetting(data.DataConfig(
		data.WithSettingName(PartMap[Part]),
		data.WithSettingEdit(Edit)))
	if err != nil {
		return nil, err
	} else {
		return PartMap[Part], nil
	}
}
