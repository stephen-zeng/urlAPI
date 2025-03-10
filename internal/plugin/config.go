package plugin

import (
	"log"
	"urlAPI/internal/data"
)

type TxtMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Txt struct {
	Model    string       `json:"model"`
	Messages []TxtMessage `json:"messages"`
}

type txtResp struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type aliImgResp struct {
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

type openaiImgResp struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

type regionResp struct {
	IPData struct {
		Info1 string `json:"info1"`
	} `json:"ipdata"`
}

// used in actual action
func fetchConfig(from string) (string, string, error) {
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{from})))
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
	IP        string
	Repo      string
}

type PluginResponse struct {
	URL          string `json:"url"`
	Response     string `json:"response"`
	InitPrompt   string `json:"init_prompt"`
	ActualPrompt string `json:"actual_prompt"`
	Context      string `json:"context"`
	Region       string `json:"region"`
	Repo         string `json:"repo"`
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
func WithIP(ip string) Option {
	return func(c *Config) {
		c.IP = ip
	}
}
func WithRepo(repo string) Option {
	return func(c *Config) {
		c.Repo = repo
	}
}
func WithSumPrompt(sumPrompt string) Option {
	return func(c *Config) {
		c.SumPrompt = sumPrompt
	}
}
func PluginConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
