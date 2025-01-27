package data

import (
	"errors"
	"math/rand"
	"time"
)

var Part = map[string][]string{
	"apiOpenai":   []string{"openai"},
	"apiAlibaba":  []string{"alibaba"},
	"apiDeekseek": []string{"deepseek"},
	"apiOtherapi": []string{"otherapi"},
	"security":    []string{"dash", "blocklist", "allow"},
	"txt":         []string{"txt", "txtrandomenabled", "txtsummaryenabled"},
	"img":         []string{"img"},
	"web":         []string{"web"},
}

func InitSetting() (string, error) {
	rand.Seed(time.Now().UnixNano())
	acsii := []int{10, 26, 26}
	acsiiPlus := []int{48, 65, 97}
	pwd := ""
	for i := 1; i <= 8; i++ {
		choose := rand.Int() % len(acsii)
		pwd += string(rand.Int()%acsii[choose] + acsiiPlus[choose])
	}
	openai := []string{"none", "https://api.openai.com/v1/chat/completions", "gpt-4o", "gpt-4o-mini", "dall-e-3", "1024x1024"}
	deepseek := []string{"none", "deepseek-chat", "deepseek-chat"}
	alibaba := []string{"none", "qwen-plus", "qwen-turbo", "wanx2.0-t2i-turbo", "1024x1024"}
	otherapi := []string{"none", "none", "none", "none"}
	dash := []string{pwd}
	blocklist := []string{}
	allow := []string{"*"}
	txt := []string{"openai", "true", "gpt-4o-mini", "gpt-4o-mini"}
	txtrandomenabled := []string{}
	txtsummaryenabled := []string{}
	img := []string{"false", "openai", "false"}
	web := []string{""}
	err := editSetting([]string{"openai", "deepseek", "alibaba", "otherapi",
		"dash", "blocklist", "allow", "txt",
		"txtrandomenabled", "txtsummaryenabled", "img", "web"},
		[][]string{openai, deepseek, alibaba, otherapi,
			dash, blocklist, allow, txt,
			txtrandomenabled, txtsummaryenabled, img, web})
	if err != nil {
		return "", errors.New("setting.restore.error")
	} else {
		return pwd, nil
	}
}

func EditSetting(part string, data [][]string) error {
	err := editSetting(Part[part], data)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func FetchSetting(data Config) ([][]string, error) {
	var ret [][]string
	var err error
	if data.Name != "" {
		ret, err = fetchSetting([]string{data.Name})
	} else if data.Part != "" {
		ret, err = fetchSetting(Part[data.Part])
	} else {
		return nil, errors.New("setting.fetch.error")
	}
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}
