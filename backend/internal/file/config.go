package file

type Config struct {
	Type string
	UUID string
	URL  string
	File string
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
func WithFile(file string) Option {
	return func(c *Config) {
		c.File = file
	}
}
func FileConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
