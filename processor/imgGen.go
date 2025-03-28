package processor

import (
	"bytes"
	"errors"
	"fmt"
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
	prompt := info.Target
	var err error
	switch info.API {
	case "alibaba":
		img, prompt, err = util.AlibabaImg(token, info.Target, info.Model, info.Size)
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
	file, err := os.Create(ImgPath + data.UUID + ".png")
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
	data.Return = fmt.Sprintf(`{"original_prompt"：%s, "actual_prompt"：%s, "url": %s}`, info.Target, prompt, info.Host+"/download?img="+data.UUID)
	data.Status = "success"
	info.Return = info.Host + "/download?img=" + data.UUID
	defer file.Close()
	return nil
}
