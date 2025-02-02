package data

import "time"

type Config struct {
	UUID   string
	Status string
	Return string
	Target string
	IP     string
	Type   string
	Name   []string
	Edit   [][]string
	Size   string
	API    string
	Token  string
	Expire time.Time
	Term   bool
	Region string
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
func WithEdit(edit [][]string) Option {
	return func(config *Config) {
		config.Edit = edit
	}
}
func WithName(name []string) Option {
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
func WithToken(token string) Option {
	return func(config *Config) {
		config.Token = token
	}
}
func WithExpire(expire time.Time) Option {
	return func(config *Config) {
		config.Expire = expire
	}
}
func WithTerm(term bool) Option {
	return func(config *Config) {
		config.Term = term
	}
}
func WithRegion(region string) Option {
	return func(config *Config) {
		config.Region = region
	}
}

func DataConfig(opts ...Option) Config {
	config := Config{
		Term: false,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
