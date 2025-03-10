package web

import (
	"golang.org/x/net/html"
	"strings"
)

type WebResponse struct {
	URL    string `json:"url"`
	Target string `json:"target"`
}

type Config struct {
	Target string `json:"target"`
	API    string `json:"api"`
}
type Option func(*Config)

func WithTarget(url string) Option {
	return func(o *Config) {
		o.Target = url
	}
}
func WithAPI(api string) Option {
	return func(o *Config) {
		o.API = api
	}
}
func WebConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

func findItem(n *html.Node) string {
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			ret += strings.TrimSpace(c.Data)
		} else if c.Type == html.ElementNode {
			ret += findItem(c)
		}
	}
	return ret
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

type ytbResp struct {
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
		Statistisc struct {
			ViewCount string `json:"viewCount"`
			LikeCount string `json:"likeCount"`
		} `json:"statistisc"`
	} `json:"items"`
}

type repoResp struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	ForksCount      float64 `json:"forks_count"`
	StargazersCount float64 `json:"stargazers_count"`
}
