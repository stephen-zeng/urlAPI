package data

import "time"

var FallbackURL = "https://raw.githubusercontent.com/stephen-zeng/urlAPI/main/fallback.png"
var Expired = 60

type Config struct {
	UUID          string
	Type          string
	API           string
	By            string
	TaskStatus    string
	TaskReturn    string
	TaskTarget    string
	TaskSize      string
	TaskRegion    string
	TaskIP        string
	TaskModel     string
	TaskDevice    string
	TaskReferer   string
	SettingName   []string
	SettingEdit   [][]string
	SessionToken  string
	SessionExpire time.Time
	SessionTerm   bool
	RepoInfo      string
}
type Option func(*Config)

func WithUUID(id string) Option {
	return func(config *Config) {
		config.UUID = id
	}
}
func WithTaskStatus(status string) Option {
	return func(config *Config) {
		config.TaskStatus = status
	}
}
func WithTaskReturn(ret string) Option {
	return func(config *Config) {
		config.TaskReturn = ret
	}
}
func WithTaskTarget(target string) Option {
	return func(config *Config) {
		config.TaskTarget = target
	}
}
func WithTaskIP(ip string) Option {
	return func(config *Config) {
		config.TaskIP = ip
	}
}
func WithType(t string) Option {
	return func(config *Config) {
		config.Type = t
	}
}
func WithSettingEdit(edit [][]string) Option {
	return func(config *Config) {
		config.SettingEdit = edit
	}
}
func WithSettingName(name []string) Option {
	return func(config *Config) {
		config.SettingName = name
	}
}
func WithTaskSize(size string) Option {
	return func(config *Config) {
		config.TaskSize = size
	}
}
func WithAPI(api string) Option {
	return func(config *Config) {
		config.API = api
	}
}
func WithSessionToken(token string) Option {
	return func(config *Config) {
		config.SessionToken = token
	}
}
func WithSessionExpire(expire time.Time) Option {
	return func(config *Config) {
		config.SessionExpire = expire
	}
}
func WithSessionTerm(term bool) Option {
	return func(config *Config) {
		config.SessionTerm = term
	}
}
func WithTaskRegion(region string) Option {
	return func(config *Config) {
		config.TaskRegion = region
	}
}
func WithTaskModel(model string) Option {
	return func(config *Config) {
		config.TaskModel = model
	}
}
func WithTaskReferer(referer string) Option {
	return func(config *Config) {
		config.TaskReferer = referer
	}
}
func WithTaskDevice(device string) Option {
	return func(config *Config) {
		config.TaskDevice = device
	}
}
func WithBy(by string) Option {
	return func(config *Config) {
		config.By = by
	}
}
func WithRepoInfo(repoInfo string) Option {
	return func(config *Config) {
		config.RepoInfo = repoInfo
	}
}
func DataConfig(opts ...Option) Config {
	config := Config{
		SessionTerm: false,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
