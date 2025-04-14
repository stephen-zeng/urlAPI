package util

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"log"
	"net"
	"net/http"
	"time"
	"urlAPI/file"
)

var (
	GlobalHTTPClient *http.Client
	font             *truetype.Font
	IPTmp            = make(map[string]string)
)

func init() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	GlobalHTTPClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	reader, err := file.Font.ReadFile("ssfonts.ttf")
	if err != nil {
		log.Println("Read font file error")
	}
	font, _ = freetype.ParseFont(reader)
	if err != nil {
		log.Println("Parse font error")
	}
}

type RegionResp struct {
	IPData struct {
		Info1 string `json:"info1"`
	} `json:"ipdata"`
}

var TypeMap = map[string]string{
	"download": "文件下载",
	"txt.gen":  "文字生成",
	"img.gen":  "图片生成",
	"rand":     "随机图片",
	"web.img":  "网页缩略图",
}

type TxtMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type TxtPayload struct {
	Model    string       `json:"model"`
	Messages []TxtMessage `json:"messages"`
}

type TxtResp struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type AlibabaImgInput struct {
	Prompt string `json:"prompt"`
}
type AlibabaImgParameters struct {
	Size string `json:"size"`
	N    int    `json:"n"`
}
type AlibabaImgPayload struct {
	Model      string               `json:"model"`
	Input      AlibabaImgInput      `json:"input"`
	Parameters AlibabaImgParameters `json:"parameters"`
}
type AlibabaImgResp struct {
	Output struct {
		TaskStatus string `json:"task_status"`
		TaskID     string `json:"task_id"`
		Results    []struct {
			OrigPrompt   string `json:"orig_prompt"`
			ActualPrompt string `json:"actual_prompt"`
			URL          string `json:"url"`
		} `json:"results"`
	} `json:"output"`
}

type OpenaiImgPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

type OpenaiImgResp struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

type BiliResp struct {
	Data struct {
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
		Stat struct {
			View     float64 `json:"view"`
			Favorite float64 `json:"favorite"`
			Like     float64 `json:"like"`
			Coin     float64 `json:"coin"`
		} `json:"stat"`
		Pic   string `json:"pic"`
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"data"`
}

type YtbResp struct {
	Items []struct {
		Snippet struct {
			Thumbnails struct {
				Standard struct {
					URL string `json:"url"`
				} `json:"standard"`
			} `json:"thumbnails"`
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
			Description  string `json:"description"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount string `json:"viewCount"`
			LikeCount string `json:"likeCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type RepoResp struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	ForksCount      float64 `json:"forks_count"`
	StargazersCount float64 `json:"stargazers_count"`
}

type RepoContentResp struct {
	DownloadURL string `json:"download_url"`
}
