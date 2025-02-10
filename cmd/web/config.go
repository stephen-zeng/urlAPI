package web

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
