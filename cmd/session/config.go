package session

import (
	"backend/internal/data"
	"time"
)

type SessionResponse struct {
	Token   string      `json:"token"`
	Name    []string    `json:"name"`
	Part    string      `json:"part"`
	Setting [][]string  `json:"setting"`
	Task    []data.Task `json:"task"`
	IP      string      `json:"ip"`
}

type Config struct {
	Type      string
	Token     string
	IP        string
	Time      time.Time
	Term      bool
	Operation string
	Part      string
	Edit      [][]string
	By        string
	Task      string
}
type Option func(*Config)

func WithType(t string) Option {
	return func(c *Config) {
		c.Type = t
	}
}
func WithToken(t string) Option {
	return func(c *Config) {
		c.Token = t
	}
}
func WithIP(ip string) Option {
	return func(c *Config) {
		c.IP = ip
	}
}
func WithTime(t time.Time) Option {
	return func(c *Config) {
		c.Time = t
	}
}
func WithTerm(term bool) Option {
	return func(c *Config) {
		c.Term = term
	}
}
func WithOperation(operation string) Option {
	return func(c *Config) {
		c.Operation = operation
	}
}
func WithPart(part string) Option {
	return func(c *Config) {
		c.Part = part
	}
}
func WithEdit(edit [][]string) Option {
	return func(c *Config) {
		c.Edit = edit
	}
}
func WithBy(by string) Option {
	return func(c *Config) {
		c.By = by
	}
}
func WithTask(task string) Option {
	return func(c *Config) {
		c.Task = task
	}
}
func SessionConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
