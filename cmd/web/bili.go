package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func biliGetStr(x float64) string {
	if x >= 10000 {
		return strconv.FormatFloat(x/10000.0, 'f', 1, 64) + "w"
	} else {
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
}

func bili(ABV, From, UUID string) (WebResponse, error) {
	var url string
	if ABV[0] == 'a' {
		ABV = ABV[2:]
		url = "https://api.bilibili.com/x/web-interface/view?aid=" + ABV
	} else if ABV[0] == 'B' {
		url = "https://api.bilibili.com/x/web-interface/view?bvid=" + ABV
	} else {
		return WebResponse{}, errors.New("Invalid ABV")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WebResponse{}, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting video info")
		return WebResponse{}, err
	}
	var info map[string]interface{}
	err = json.Unmarshal(jsonResp, &info)
	if err != nil {
		return WebResponse{}, err
	}
	info = info["data"].(map[string]interface{})
	picURL := info["pic"].(string)
	name := info["title"].(string)
	author := info["owner"].(map[string]interface{})["name"].(string)
	description := info["desc"].(string)
	info = info["stat"].(map[string]interface{})
	view := biliGetStr(info["view"].(float64))
	favorite := biliGetStr(info["favorite"].(float64))
	like := biliGetStr(info["like"].(float64))
	coin := biliGetStr(info["coin"].(float64))
	err = drawVideo(picURL, name, author, description, view, favorite, like, coin, UUID)
	if err != nil {
		log.Println("Error when drawing the img")
		return WebResponse{}, err
	}
	return WebResponse{
		URL: From + "/download?img=" + UUID,
	}, err
}
