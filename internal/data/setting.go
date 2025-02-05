package data

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func InitSetting(data Config) (string, error) {
	if data.Type != "restore" && db.Migrator().HasTable(&Setting{}) {
		return "", nil
	}
	db.AutoMigrate(&Setting{})
	rand.Seed(time.Now().UnixNano())
	acsii := []int{10, 26, 26}
	acsiiPlus := []int{48, 65, 97}
	pwd := ""
	for i := 1; i <= 8; i++ {
		choose := rand.Int() % len(acsii)
		pwd += string(rand.Int()%acsii[choose] + acsiiPlus[choose])
	}
	openai := []string{"", "gpt-4o", "gpt-4o-mini", "dall-e-3", "1024x1024", "https://api.openai.com/v1/chat/completions"}
	deepseek := []string{"", "deepseek-chat", "deepseek-chat"}
	alibaba := []string{"", "qwen-plus", "qwen-turbo", "wanx2.0-t2i-turbo", "1024x768"}
	otherapi := []string{"", "", "", ""}

	dash := []string{fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))}
	dashallowedip := []string{"*"}
	allowedref := []string{"*"}

	txt := []string{"false", "alibaba", "alibaba", "60"}
	txtgenenabled := []string{"_"}
	txtsumenabled := []string{"_"}

	img := []string{"false", "alibaba", "60", "https://raw.githubusercontent.com/stephen-zeng/urlAPI/img/master/fallback.png"}

	web := []string{"false", "false", "alibaba"}
	webimgallowed := []string{"_"}
	websumblocked := []string{"_"}

	rd := []string{"false", "https://raw.githubusercontent.com", "https://raw.githubusercontent.com/stephen-zeng/urlAPI/master/fallback.png"}
	err := editSetting(
		[]string{"openai", "deepseek", "alibaba", "otherapi", "dash", "dashallowedip", "allowedref", "txt", "txtgenenabled", "txtsumenabled", "img", "web", "webimgallowed", "websumblocked", "rand"},
		[][]string{openai, deepseek, alibaba, otherapi, dash, dashallowedip, allowedref, txt, txtgenenabled, txtsumenabled, img, web, webimgallowed, websumblocked, rd})
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Println("Initialized Setting")
		return pwd, nil
	}
}

func EditSetting(data Config) error {
	err := editSetting(data.SettingName, data.SettingEdit)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func FetchSetting(data Config) ([][]string, error) {
	ret, err := fetchSetting(data.SettingName)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return ret, nil
	}
}
