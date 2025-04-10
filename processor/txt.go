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
		data.Model = info.Model
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
	img, err := util.DrawTxt(response)
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor Txt"), err)
	}
	file, err := os.Create(ImgPath + data.UUID + ".png")
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
	
	data.Return = fmt.Sprintf(`{"prompt": "%s", "response": "%s", "url": "%s"}`, info.Target, response, info.Host+"/download?img="+data.UUID)
	data.Status = "success"
	info.Return = info.Host + "/download?img=" + data.UUID
	defer file.Close()
	return nil
}
