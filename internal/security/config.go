package security

type Config struct {
	IP     string
	Pwd    string
	Type   string
	API    string
	Repo   string
	User   string
	Domain string
	Target string
}
type Option func(*Config)

func WithIP(ip string) Option {
	return func(config *Config) {
		config.IP = ip
	}
}
func WithPwd(pwd string) Option {
	return func(config *Config) {
		config.Pwd = pwd
	}
}
func WithType(t string) Option {
	return func(config *Config) {
		config.Type = t
	}
}
func WithAPI(api string) Option {
	return func(config *Config) {
		config.API = api
	}
}
func WithRepo(repo string) Option {
	return func(config *Config) {
		config.Repo = repo
	}
}
func WithUser(user string) Option {
	return func(config *Config) {
		config.User = user
	}
}
func WithDomain(domain string) Option {
	return func(config *Config) {
		config.Domain = domain
	}
}
func WithTarget(target string) Option {
	return func(config *Config) {
		config.Target = target
	}
}
func SecurityConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
