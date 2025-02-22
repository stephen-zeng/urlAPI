package set

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"time"
	"urlAPI/internal/data"
)

func RePwd() (SetResponse, error) {
	dat, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"dash"})))
	if err != nil {
		return SetResponse{}, err
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
		data.WithSettingName([]string{"dash"}),
		data.WithSettingEdit(dat)))
	return SetResponse{
		Pwd: pwd,
	}, err
}

func Restore() (SetResponse, error) {
	err := data.InitSession(data.DataConfig(data.WithType("restore")))
	if err != nil {
		return SetResponse{}, err
	}
	pwd, err := data.InitSetting(data.DataConfig(data.WithType("restore")))
	return SetResponse{
		Pwd: pwd,
	}, err
}

func Clear() error {
	os.RemoveAll("assets/img")
	os.Mkdir("assets/img", 0777)
	err := data.InitTask(data.DataConfig(data.WithType("restore")))
	return err
}

func ClearIP() error {
	err := data.EditSetting(data.DataConfig(
		data.WithSettingName([]string{"dashallowedip"}),
		data.WithSettingEdit([][]string{{"*"}})))
	return err
}

func Update() error {
	_, err := data.InitSetting(data.DataConfig(data.WithType("update")))
	return err
}
