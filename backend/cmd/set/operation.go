package set

import (
	"backend/internal/data"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

var Part = map[string][]string{
	"openai":   []string{"openai"},
	"deepseek": []string{"deepseek"},
	"alibaba":  []string{"alibaba"},
	"otherapi": []string{"otherapi"},
	"security": []string{"dash", "dashallowedip", "allowedref"},
	"txt":      []string{"txt", "txtgenenabled", "txtsumenabled"},
	"img":      []string{"img"},
	"web":      []string{"web", "webimgallowed", "websumblocked"},
}

func fetch(dat Config) ([][]string, error) {
	list, err := data.FetchSetting(
		data.DataConfig(
			data.WithName(Part[dat.part])))
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func edit(dat Config) error {
	return data.EditSetting(data.DataConfig(
		data.WithName(Part[dat.part]),
		data.WithEdit(dat.edit)))
}

func repwd() (string, error) {
	dat, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"dash"})))
	if err != nil {
		return "", err
	}
	rand.Seed(time.Now().UnixNano())
	acsii := []int{10, 26, 26}
	acsiiPlus := []int{48, 65, 97}
	pwd := ""
	for i := 1; i <= 8; i++ {
		choose := rand.Int() % len(acsii)
		pwd += string(rand.Int()%acsii[choose] + acsiiPlus[choose])
	}
	dat[0][0] = fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))
	err = data.EditSetting(data.DataConfig(
		data.WithName([]string{"dash"}),
		data.WithEdit(dat)))
	return pwd, err
}
