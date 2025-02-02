package session

import "time"

type Config struct {
	Type  string
	Token string
	Json  map[string]interface{}
	IP    string
	Time  time.Time
	Term  bool
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
func WithJson(j map[string]interface{}) Option {
	return func(c *Config) {
		c.Json = j
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
func WithTerm() Option {
	return func(c *Config) {
		c.Term = true
	}
}
func SessionConfig(opts ...Option) Config {
	config := Config{
		Term: false,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
