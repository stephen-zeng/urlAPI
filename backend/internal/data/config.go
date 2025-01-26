package data

type Config struct {
	UUID   string
	Status string
	Return string
	Target string
	IP     string
	Type   string
}
type Option func(*Config)

func WithUUID(id string) Option {
	return func(config *Config) {
		config.UUID = id
	}
}
func WithStatus(status string) Option {
	return func(config *Config) {
		config.Status = status
	}
}
func WithReturn(ret string) Option {
	return func(config *Config) {
		config.Return = ret
	}
}
func WithTarget(target string) Option {
	return func(config *Config) {
		config.Target = target
	}
}
func WithIP(ip string) Option {
	return func(config *Config) {
		config.IP = ip
	}
}
func WithType(t string) Option {
	return func(config *Config) {
		config.Type = t
	}
}

func DataConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
