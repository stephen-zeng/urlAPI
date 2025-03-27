package processor

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/url"
	"os"
	"urlAPI/database"
	"urlAPI/util"
)

func getBiliABV(URL string) string {
	for i := 31; i < len(URL); i++ {
		if URL[i] == '/' || URL[i] == '?' {
			return URL[31:i]
		}
	}
	return ""
}

func getYtbID(URL string) string {
	for i := 32; i < len(URL); i++ {
		if URL[i] == '&' {
			return URL[32:i]
		}
	}
	return URL[32:]
}

func (info *WebImg) Process(data *database.Task) error {
	var img []byte
	info.Return = database.SettingMap["web"][4]
	urlParse, err := url.Parse(info.Return)
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor WebImg"), err)
	}
	if _, ok := WebImgMap[urlParse.Host]; !ok {
		data.Status = "failed"
		data.Return = "Invalid URL"
		return errors.New("Processor WebImg Invalid URL")
	}
	info.API = WebImgMap[urlParse.Host]
	data.API = info.API
	switch info.API {
	case "bilibili":
		img, err = util.Bili(getBiliABV(info.Target))
	case "ytb":
		if len(database.SettingMap["web"]) < 7 {
			err = errors.Join(errors.New("Processor WebImg Not valid ytb token"))
			break
		}
		token := database.SettingMap["web"][6]
		img, err = util.Ytb(getYtbID(info.Target), token)
	case "arxiv":
		img, err = util.Arxiv(info.Target)
	case "ITHome":
		api := database.SettingMap["txt"][1]
		token := database.SettingMap[api][0]
		model := database.SettingMap[api][2]
		context := database.SettingMap["context"][1]
		endpoint := getEndpoint(api)
		img, err = util.ITHome(info.Target, endpoint, token, model, context)
	case "Repo":
		token := ""
		if len(database.SettingMap["token"]) > 5 {
			token = database.SettingMap["token"][5]
		}
		img, err = util.Repo(info.Target, token)
	}
	id := uuid.New().String()
	file, err := os.Create(ImgPath + id + ".png")
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor WebImg"), err)
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(img))
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.Join(errors.New("Processor WebImg"), err)
	}
	data.Status = "success"
	data.Return = id
	info.Return = info.Host + "/download?img=" + id + ".png"
	return nil
}
