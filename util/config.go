package util

import (
	"net"
	"net/http"
	"time"
)

type RegionResp struct {
	IPData struct {
		Info1 string `json:"info1"`
	} `json:"ipdata"`
}

var TypeMap = map[string]string{
	"download": "文件下载",
	"txt.gen":  "文字生成",
	"txt.sum":  "文字总结",
	"img.gen":  "图片生成",
	"rand":     "随机图片",
	"web.img":  "网页缩略图",
}

var GlobalHTTPClient *http.Client

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
}
