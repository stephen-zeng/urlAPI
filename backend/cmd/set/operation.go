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

func fetch(part string) (SetResponse, error) {
	list, err := data.FetchSetting(
		data.DataConfig(
			data.WithName(Part[part])))
	if err != nil {
		return SetResponse{}, err
	} else {
		return SetResponse{
			Name:    Part[part],
			Setting: list,
		}, nil
	}
}

func edit(part string, edit [][]string) (SetResponse, error) {
	err := data.EditSetting(data.DataConfig(
		data.WithName(Part[part]),
		data.WithEdit(edit)))
	if err != nil {
		return SetResponse{}, err
	} else {
		return SetResponse{
			Name: Part[part],
		}, err
	}
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
