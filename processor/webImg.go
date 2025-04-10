package processor

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
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
	return URL[31:]
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
	urlParse, err := url.Parse(info.Target)
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.WithStack(err)
	}
	info.API = urlParse.Host
	data.API = info.API
	switch info.API {
	case "www.bilibili.com":
		img, err = util.Bili(getBiliABV(info.Target))
	case "www.youtube.com":
		token := database.SettingMap["web"][6]
		img, err = util.Ytb(getYtbID(info.Target), token)
	case "arxiv.org":
		img, err = util.Arxiv(info.Target)
	case "www.ithome.com":
		api := database.SettingMap["txt"][1]
		token := database.SettingMap[api][0]
		model := database.SettingMap[api][2]
		context := database.SettingMap["context"][1]
		endpoint := getEndpoint(api)
		img, err = util.ITHome(info.Target, endpoint, token, model, context)
	case "github.com", "gitee.com":
		token := database.SettingMap["web"][5]
		img, err = util.Repo(info.Target, token)
	default:
		data.Status = "failed"
		data.Return = "Invalid URL"
		return errors.WithStack(errors.New("Invalid URL"))
	}
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.WithStack(err)
	}
	file, _ := os.Create(ImgPath + data.UUID + ".png")
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(img))
	if err != nil {
		data.Status = "failed"
		data.Return = err.Error()
		return errors.WithStack(err)
	}
	data.Status = "success"
	data.Return = fmt.Sprintf(`{"url": "%s"}`, info.Host+"/download?img="+data.UUID)
	info.Return = info.Host + "/download?img=" + data.UUID
	return nil
}
