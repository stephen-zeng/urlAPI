package processor

import (
	"github.com/pkg/errors"
	"os"
	"urlAPI/database"
)

func (info *Download) Process(data *database.Task) error {
	img, err := os.ReadFile(ImgPath + info.Target + ".png")
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		info.ReturnError = database.SettingMap["web"][4]
		return errors.WithStack(err)
	}
	data.Return = ImgPath + info.Target + ".png"
	data.Status = "success"
	info.Return = img
	return nil
}
