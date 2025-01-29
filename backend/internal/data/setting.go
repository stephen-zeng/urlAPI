package data

import (
	"log"
	"math/rand"
	"time"
)

var Part = map[string][]string{
	"apiOpenai":   []string{"openai"},
	"apiAlibaba":  []string{"alibaba"},
	"apiDeekseek": []string{"deepseek"},
	"apiOtherapi": []string{"otherapi"},
	"security":    []string{"dash", "blocklist", "allow"},
	"txt":         []string{"txt", "txtgenenabled", "txtsumenabled"},
	"img":         []string{"img"},
	"web":         []string{"web"},
}

func InitSetting(data Config) (string, error) {
	if data.Type != "restore" && db.Migrator().HasTable(&Setting{}) {
		return "", nil
	}
	log.Println("Start init setting")
	rand.Seed(time.Now().UnixNano())
	acsii := []int{10, 26, 26}
	acsiiPlus := []int{48, 65, 97}
	pwd := ""
	for i := 1; i <= 8; i++ {
		choose := rand.Int() % len(acsii)
		pwd += string(rand.Int()%acsii[choose] + acsiiPlus[choose])
	}
	openai := []string{"none", "gpt-4o", "gpt-4o-mini", "dall-e-3", "1024x1024", "https://api.openai.com/v1/chat/completions"}
	deepseek := []string{"none", "deepseek-chat", "deepseek-chat"}
	alibaba := []string{"none", "qwen-plus", "qwen-turbo", "wanx2.0-t2i-turbo", "1024x1024"}
	otherapi := []string{"none", "none", "none", "none"}
	dash := []string{pwd, "*"}
	blocklist := []string{}
	allow := []string{"*"}
	txt := []string{"openai", "true", "gpt-4o-mini", "true", "gpt-4o-mini"}
	txtgenenabled := []string{"laugh", "poem", "sentence", "other"}
	txtsumenabled := []string{}
	img := []string{"false", "openai", "false"}
	web := []string{""}
	rand := []string{"https://gh.qwqwq.com.cn"}
	err := editSetting([]string{"openai", "deepseek", "alibaba", "otherapi",
		"dash", "blocklist", "allow", "txt",
		"txtgenenabled", "txtsumenabled", "img", "web", "rand"},
		[][]string{openai, deepseek, alibaba, otherapi,
			dash, blocklist, allow, txt,
			txtgenenabled, txtsumenabled, img, web, rand})
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		return pwd, nil
	}
}

func EditSetting(part string, data [][]string) error {
	err := editSetting(Part[part], data)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return nil, err
	}
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return ret, nil
	}
}
