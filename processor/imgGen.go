package processor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"urlAPI/database"
	"urlAPI/util"
)

func (info *ImgGen) Process(data *database.Task) error {
	info.Return = database.SettingMap["img"][3]
	if info.API == "" {
		info.API = database.SettingMap["img"][1]
		data.API = info.API
	}
	if info.Model == "" {
		info.Model = database.SettingMap[info.API][3]
		data.Model = info.Model
	}
	if info.Size == "" {
		info.Size = database.SettingMap[info.API][4]
		data.Size = info.Size
	}
	token := database.SettingMap[info.API][0]
	var img []byte
	var prompt string
	var err error
	switch info.API {
	case "alibaba":
		img, prompt, err = util.AlibabaImg(token, info.Target, info.Model, info.Size)
		prompt = fmt.Sprintf(`原始Prompt："%s"，实际Prompt："%s"。`, info.Target, prompt)
	case "openai":
		prompt = fmt.Sprintf(`Prompt："%s"。`, info.Target)
		img, err = util.OpenaiImg(database.SettingMap["openai"][5],
			token, info.Target, info.Model, info.Size)
	default:
		data.Status = "failed"
		data.Return = "Imggen Process invalid API"
		return errors.New("Imggen Process invalid API")
	}
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.New("Imggen Process " + err.Error())
	}
	id := uuid.New().String()
	file, err := os.Create(ImgPath + id + ".png")
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.New("Imggen Process " + err.Error())
	}
	_, err = io.Copy(file, bytes.NewReader(img))
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.New("Imggen Process " + err.Error())
	}
	data.Return = prompt + "UUID: " + id
	data.Status = "success"
	info.Return = info.Host + "/download?img=" + id
	defer file.Close()
	return nil
}
