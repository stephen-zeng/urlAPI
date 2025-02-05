package set

import (
	"backend/internal/data"
	"os"
)

func Fetch(dat Config) (SetResponse, error) {
	return fetch(dat.part)
}

func Edit(dat Config) (SetResponse, error) {
	return edit(dat.part, dat.edit)
}

func RePwd() (SetResponse, error) {
	pwd, err := repwd()
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
		data.WithName([]string{"dashallowedip"}),
		data.WithEdit([][]string{{"*"}})))
	return err
}
