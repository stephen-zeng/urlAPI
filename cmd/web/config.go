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
