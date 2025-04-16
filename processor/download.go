package processor

import (
	"github.com/pkg/errors"
	"os"
	"urlAPI/database"
	"urlAPI/file"
)

func (info *Download) Process(data *database.Task) error {
	var img []byte
	var err error
	switch info.Target {
	case "empty":
		img, err = file.EmptyPNG.ReadFile("empty.png")
	default:
		img, err = os.ReadFile(ImgPath + info.Target + ".png")
	}
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
