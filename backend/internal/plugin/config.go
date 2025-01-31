package plugin

import (
	"backend/internal/data"
	"log"
)

type TxtMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Txt struct {
	Model    string       `json:"model"`
	Messages []TxtMessage `json:"messages"`
}

func fetchConfig(from string) (string, string, error) {
	config, err := data.FetchSetting(data.DataConfig(data.WithName([]string{from})))
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	token := config[0][0]
	url := config[0][len(config[0])-1]
	return url, token, nil
}

type Config struct {
	API       string
	GenPrompt string
	SumPrompt string
	Size      string
	ImgPrompt string
	Model     string
	User      string
	Repo      string
}

type Option func(*Config)

func WithAPI(api string) Option {
	return func(c *Config) {
		c.API = api
	}
}
func WithGenPrompt(genPrompt string) Option {
	return func(c *Config) {
		c.GenPrompt = genPrompt
	}
}
func WithSumPrompt(sumPrompt string) Option {
	return func(c *Config) {
		c.SumPrompt = sumPrompt
	}
}
func WithSize(size string) Option {
	return func(c *Config) {
		c.Size = size
	}
}
func WithImgPrompt(imgPrompt string) Option {
	return func(c *Config) {
		c.ImgPrompt = imgPrompt
	}
}
func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}
func WithUser(user string) Option {
	return func(c *Config) {
		c.User = user
	}
}
func WithRepo(repo string) Option {
	return func(c *Config) {
		c.Repo = repo
	}
}

func PluginConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
