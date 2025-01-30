package data

type Config struct {
	UUID   string
	Status string
	Return string
	Target string
	IP     string
	Type   string
	Part   string
	Name   string
	Size   string
	API    string
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
func WithPart(part string) Option {
	return func(config *Config) {
		config.Part = part
	}
}
func WithName(name string) Option {
	return func(config *Config) {
		config.Name = name
	}
}
func WithSize(size string) Option {
	return func(config *Config) {
		config.Size = size
	}
}
func WithAPI(api string) Option {
	return func(config *Config) {
		config.API = api
	}
}

func DataConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
