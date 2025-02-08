package file

import "image"

type Config struct {
	Type string
	UUID string
	URL  string
	Img  image.Image
}
type Option func(*Config)

func WithType(t string) Option {
	return func(c *Config) {
		c.Type = t
	}
}
func WithUUID(uuid string) Option {
	return func(c *Config) {
		c.UUID = uuid
	}
}
func WithURL(url string) Option {
	return func(c *Config) {
		c.URL = url
	}
}
func WithImg(img image.Image) Option {
	return func(c *Config) {
		c.Img = img
	}
}
func FileConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
