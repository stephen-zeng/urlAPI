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

	dash := []string{fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))}
	dashallowedip := []string{"*"}
	allowedref := []string{"*"}

	txt := []string{"openai", "true", "openai"}
	txtgenenabled := []string{"laugh", "poem", "sentence", "other"}
	txtsumenabled := []string{"pdf", "word"}

	img := []string{"false", "openai", "https://gh.qwqwq.com.cn"}

	web := []string{"true", "true", "gpt-4o-mini"}
	webimgallowed := []string{""}
	websumblocked := []string{""}
	err := editSetting([]string{"openai", "deepseek", "alibaba", "otherapi", "dash", "dashallowedip", "allowedref", "txt", "txtgenenabled", "txtsumenabled", "img", "web", "webimgallowed", "websumblocked"},
		[][]string{openai, deepseek, alibaba, otherapi, dash, dashallowedip, allowedref, txt, txtgenenabled, txtsumenabled, img, web, webimgallowed, websumblocked})
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Println("Initialized Setting")
		return pwd, nil
	}
}

func EditSetting(data Config) error {
	err := editSetting(data.Name, data.Edit)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func FetchSetting(data Config) ([][]string, error) {
	ret, err := fetchSetting(data.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return ret, nil
	}
}
