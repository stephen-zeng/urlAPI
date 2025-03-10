package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"urlAPI/internal/client"
)

func ytb(ID, From, UUID, Token string) (WebResponse, error) {
	url := "https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics&id=" + ID + "&key=" + Token
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
	var info ytbResp
	err = json.Unmarshal(jsonResp, &info)
	name := info.Items[0].Snippet.Title
	author := info.Items[0].Snippet.ChannelTitle
	description := info.Items[0].Snippet.Description
	picURL := info.Items[0].Snippet.Thumbnails.Standard.URL
	view := info.Items[0].Statistisc.ViewCount
	like := info.Items[0].Statistisc.LikeCount
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
