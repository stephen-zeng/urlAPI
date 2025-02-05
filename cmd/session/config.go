package session

import (
	"backend/internal/data"
	"time"
)

type SessionResponse struct {
	SessionToken string              `json:"session_token"`
	SessionIP    string              `json:"session_ip"`
	SettingName  []string            `json:"setting_name"`
	SettingPart  string              `json:"setting_part"`
	SettingData  [][]string          `json:"setting_data"`
	TaskData     []data.Task         `json:"task_data"`
	RepoData     []data.RepoResponse `json:"repo_data"`
}

type Config struct {
	Operation    string
	SessionType  string
	SessionToken string
	SessionIP    string
	SessionTime  time.Time
	SessionTerm  bool
	SettingPart  string
	SettingEdit  [][]string
	TaskCatagory string
	TaskBy       string
	RepoAPI      string
	RepoInfo     string
	RepoUUID     string
}
type Option func(*Config)

func WithSessionType(t string) Option {
	return func(c *Config) {
		c.SessionType = t
	}
}
func WithSessionToken(t string) Option {
	return func(c *Config) {
		c.SessionToken = t
	}
}
func WithSessionIP(ip string) Option {
	return func(c *Config) {
		c.SessionIP = ip
	}
}
func WithSessionTime(t time.Time) Option {
	return func(c *Config) {
		c.SessionTime = t
	}
}
func WithSessionTerm(term bool) Option {
	return func(c *Config) {
		c.SessionTerm = term
	}
}
func WithOperation(operation string) Option {
	return func(c *Config) {
		c.Operation = operation
	}
}
func WithSettingPart(part string) Option {
	return func(c *Config) {
		c.SettingPart = part
	}
}
func WithSettingEdit(edit [][]string) Option {
	return func(c *Config) {
		c.SettingEdit = edit
	}
}
func WithTaskCatagory(by string) Option {
	return func(c *Config) {
		c.TaskCatagory = by
	}
}
func WithTaskBy(task string) Option {
	return func(c *Config) {
		c.TaskBy = task
	}
}
func WithRepoAPI(api string) Option {
	return func(c *Config) {
		c.RepoAPI = api
	}
}
func WithRepoInfo(info string) Option {
	return func(c *Config) {
		c.RepoInfo = info
	}
}
func WithRepoUUID(uuid string) Option {
	return func(c *Config) {
		c.RepoUUID = uuid
	}
}
func SessionConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
