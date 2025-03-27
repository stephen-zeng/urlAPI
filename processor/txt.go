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

func (info *TxtGen) Process(data *database.Task) error {
	info.Return = database.SettingMap["txt"][4]
	if _, ok := database.PromptMap[info.Target]; ok {
		info.Target = database.SettingMap["prompt"][database.PromptMap[info.Target]]
	}
	if info.API == "" {
		info.API = database.SettingMap["txt"][1]
		data.API = info.API
	}
	if info.Model == "" {
		info.Model = database.SettingMap[info.API][1]
	}
	token := database.SettingMap[info.API][0]
	context := database.SettingMap["context"][0]
	endpoint := getEndpoint(info.API)
	if endpoint == "" {
		data.Status = "failed"
		data.Return = "Unknown API"
		return errors.Join(errors.New("Processor Txt Unknown API"))
	}
	response, err := util.Txt(endpoint, token, info.Model, context, info.Target)
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor Txt"), err)
	}
	if info.Type == "json" {
		data.Status = "success"
		data.Return = response
		info.Return = data.Return
		return nil
	}
	img, err := util.DrawTxt(response)
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor Txt"), err)
	}
	id := uuid.New().String()
	file, err := os.Create(ImgPath + id + ".png")
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor Txt"), err)
	}
	_, err = io.Copy(file, bytes.NewReader(img))
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor Txt"), err)
	}
	data.Return = fmt.Sprintf(`Prompt: %s, UUID: %s`, info.Target, id)
	data.Status = "success"
	info.Return = info.Host + "/download?img=" + id
	defer file.Close()
	return nil
}
