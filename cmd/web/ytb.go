package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

func ytb(ID, From, UUID, Token string) (WebResponse, error) {
	url := "https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics&id=" + ID + "&key=" + Token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WebResponse{}, err
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return WebResponse{}, err
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Error getting video info")
		return WebResponse{}, errors.Join(err, errors.New(resp.Status))
	}
	var info map[string]interface{}
	err = json.Unmarshal(jsonResp, &info)
	if err != nil {
		return WebResponse{}, err
	}
	info = info["items"].([]interface{})[0].(map[string]interface{})
	snippet := info["snippet"].(map[string]interface{})
	statistics := info["statistics"].(map[string]interface{})
	name := snippet["title"].(string)
	author := snippet["channelTitle"].(string)
	description := snippet["description"].(string)
	picURL := snippet["thumbnails"].(map[string]interface{})["standard"].(map[string]interface{})["url"].(string)
	view := statistics["viewCount"].(string)
	like := statistics["likeCount"].(string)
	favorite := "N/A"
	coin := "N/A"
	err = drawVideo(picURL, name, author, description, view, favorite, like, coin, UUID)
	if err != nil {
		log.Println("Error when drawing the img")
		return WebResponse{}, err
	}
	return WebResponse{
		URL: From + "/download?img=" + UUID,
	}, err

}
