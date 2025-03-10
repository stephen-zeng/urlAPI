package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"urlAPI/internal/client"
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
	resp, err := client.GlobalHTTPClient.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting video info")
		return WebResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	var info BiliResp
	err = json.Unmarshal(jsonResp, &info)
	if err != nil {
		return WebResponse{}, err
	}
	picURL := info.Data.Pic
	name := info.Data.Title
	author := info.Data.Owner.Name
	description := info.Data.Desc
	view := biliGetStr(info.Data.Stat.View)
	favorite := biliGetStr(info.Data.Stat.Favorite)
	like := biliGetStr(info.Data.Stat.Like)
	coin := biliGetStr(info.Data.Stat.Coin)
	err = drawVideo(picURL, name, author, description, view, favorite, like, coin, UUID)
	if err != nil {
		log.Println("Error when drawing the img")
		return WebResponse{}, err
	}
	return WebResponse{
		URL: From + "/download?img=" + UUID,
	}, err
}
