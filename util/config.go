package util

import (
	"log"
	"net"
	"net/http"
	"time"
	"urlAPI/file"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

/**
 * @file config.go
 * @brief 全局客户端、字体与外部 API 数据结构定义。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

var (
	GlobalHTTPClient *http.Client
	font             *truetype.Font
	IPTmp            = make(map[string]string)
)

/**
 * @brief 初始化全局 HTTP 客户端和字体资源。
 */
func init() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
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
		return
	}
	font, err = freetype.ParseFont(reader)
	if err != nil {
		log.Println("Parse font error")
	}
}

/**
 * @brief IP 地理位置接口响应结构。
 */
type RegionResp struct {
	Data struct {
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
	} `json:"data"`
}

/** @brief 业务类型到中文描述的映射表。 */
var TypeMap = map[string]string{
	"download": "文件下载",
	"txt.gen":  "文字生成",
	"img.gen":  "图片生成",
	"rand":     "随机图片",
	"web.img":  "网页缩略图",
}

/**
 * @brief 文本对话消息结构。
 */
type TxtMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

/**
 * @brief 文本生成请求载荷。
 */
type TxtPayload struct {
	Model            string       `json:"model"`
	Messages         []TxtMessage `json:"messages"`
	Temperature      float64      `json:"temperature,omitempty"`
	MaxTokens        int          `json:"max_tokens,omitempty"`
	TopP             float64      `json:"top_p,omitempty"`
	PresencePenalty  float64      `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64      `json:"frequency_penalty,omitempty"`
}

/**
 * @brief 文本生成接口响应结构。
 */
type TxtResp struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

/** @brief 阿里云文生图输入结构。 */
type AlibabaImgInput struct {
	Prompt string `json:"prompt"`
}

/** @brief 阿里云文生图参数结构。 */
type AlibabaImgParameters struct {
	Size string `json:"size"`
	N    int    `json:"n"`
}

/** @brief 阿里云文生图请求载荷。 */
type AlibabaImgPayload struct {
	Model      string               `json:"model"`
	Input      AlibabaImgInput      `json:"input"`
	Parameters AlibabaImgParameters `json:"parameters"`
}

/** @brief 阿里云文生图响应结构。 */
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

/** @brief OpenAI 图像生成请求载荷。 */
type OpenaiImgPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

/** @brief OpenAI 图像生成响应结构。 */
type OpenaiImgResp struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

/** @brief Bilibili 视频接口响应结构。 */
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

/** @brief YouTube 视频接口响应结构。 */
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

/** @brief 代码仓库信息接口响应结构。 */
type RepoResp struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	ForksCount      float64 `json:"forks_count"`
	StargazersCount float64 `json:"stargazers_count"`
}

/** @brief 仓库内容列表接口响应结构。 */
type RepoContentResp struct {
	DownloadURL string `json:"download_url"`
}
