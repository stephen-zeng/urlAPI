package set

type SetResponse struct {
	Name    []string   `json:"name"`
	Setting [][]string `json:"setting"`
	Pwd     string     `json:"pwd"`
}
type Config struct {
	part string
	edit [][]string
}
type Option func(*Config)

func WithPart(part string) Option {
	return func(c *Config) {
		c.part = part
	}
}
func WithEdit(edit [][]string) Option {
	return func(c *Config) {
		c.edit = edit
	}
}
func SetConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
