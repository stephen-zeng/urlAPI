// Package trust holds the HTTP trust configuration: which peers may set
// forwarding headers, and which origin generated URLs should use.
package trust

import (
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"sync/atomic"

	"github.com/pkg/errors"
)

const (
	// EnvTrustedProxies lists IPs/CIDRs whose X-Forwarded-For, X-Real-IP and
	// X-Forwarded-Proto headers are honoured. "none" trusts no proxy.
	EnvTrustedProxies = "URLAPI_TRUSTED_PROXIES"
	// EnvPublicURL sets the external origin (e.g. https://api.example.com)
	// used for absolute URLs in JSON responses instead of the Host header.
	EnvPublicURL = "URLAPI_PUBLIC_URL"
	// EnvDashboardOrigins lists extra origins allowed to call /session
	// cross-origin (for example a separately hosted dashboard).
	EnvDashboardOrigins = "URLAPI_DASHBOARD_ORIGINS"
)

// DefaultTrustedProxies covers reverse proxies on the same host or a private
// network (Docker bridge, LAN), while ignoring forwarding headers sent
// directly by public clients.
var DefaultTrustedProxies = []string{
	"127.0.0.0/8", "::1/128",
	"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7",
}

type Config struct {
	Proxies          []netip.Prefix
	PublicOrigin     string
	DashboardOrigins []string
}

var current atomic.Pointer[Config]

func init() {
	cfg, _ := Parse("", "", "")
	Set(cfg)
}

// Set installs cfg as the process-wide trust configuration.
func Set(cfg Config) { current.Store(&cfg) }

// Get returns the process-wide trust configuration.
func Get() Config { return *current.Load() }

// Load reads the configuration from the environment.
func Load() (Config, error) {
	return Parse(os.Getenv(EnvTrustedProxies), os.Getenv(EnvPublicURL), os.Getenv(EnvDashboardOrigins))
}

// Parse builds a Config. An empty proxies string selects the defaults.
func Parse(proxies, publicURL, dashboardOrigins string) (Config, error) {
	var cfg Config
	var entries []string
	switch strings.ToLower(strings.TrimSpace(proxies)) {
	case "":
		entries = DefaultTrustedProxies
	case "none":
	default:
		entries = splitList(proxies)
	}
	for _, entry := range entries {
		prefix, err := parsePrefix(entry)
		if err != nil {
			return Config{}, errors.Wrapf(err, "%s: invalid entry %q", EnvTrustedProxies, entry)
		}
		cfg.Proxies = append(cfg.Proxies, prefix)
	}
	if strings.TrimSpace(publicURL) != "" {
		origin, err := normalizeOrigin(publicURL)
		if err != nil {
			return Config{}, errors.Wrap(err, EnvPublicURL)
		}
		cfg.PublicOrigin = origin
	}
	for _, origin := range splitList(dashboardOrigins) {
		normalized, err := normalizeOrigin(origin)
		if err != nil {
			return Config{}, errors.Wrap(err, EnvDashboardOrigins)
		}
		cfg.DashboardOrigins = append(cfg.DashboardOrigins, normalized)
	}
	return cfg, nil
}

func splitList(s string) []string {
	var ret []string
	for _, item := range strings.Split(s, ",") {
		if item = strings.TrimSpace(item); item != "" {
			ret = append(ret, item)
		}
	}
	return ret
}

func parsePrefix(s string) (netip.Prefix, error) {
	if strings.Contains(s, "/") {
		prefix, err := netip.ParsePrefix(s)
		return prefix.Masked(), err
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Prefix{}, err
	}
	return netip.PrefixFrom(addr, addr.BitLen()), nil
}

func normalizeOrigin(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil ||
		(u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return "", errors.Errorf("%q must be an origin such as https://api.example.com", raw)
	}
	return strings.ToLower(u.Scheme) + "://" + strings.ToLower(u.Host), nil
}

// ProxyStrings renders the trusted proxies for gin.Engine.SetTrustedProxies.
func (c Config) ProxyStrings() []string {
	ret := make([]string, 0, len(c.Proxies))
	for _, prefix := range c.Proxies {
		ret = append(ret, prefix.String())
	}
	return ret
}

// IsTrustedPeer reports whether the direct TCP peer of r is a trusted proxy.
func (c Config) IsTrustedPeer(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	for _, prefix := range c.Proxies {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

// Scheme returns "https" or "http" for the original client request. The
// X-Forwarded-Proto header is honoured only from trusted proxies and only
// for those two values.
func (c Config) Scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if c.IsTrustedPeer(r) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if i := strings.IndexByte(proto, ','); i >= 0 {
			proto = proto[:i]
		}
		switch strings.ToLower(strings.TrimSpace(proto)) {
		case "https":
			return "https"
		case "http":
			return "http"
		}
	}
	return "http"
}

// Origin returns the external origin for r: the configured public origin
// when set, otherwise one derived from the scheme and a validated Host
// header. It returns "" when no trustworthy origin can be built.
func (c Config) Origin(r *http.Request) string {
	if c.PublicOrigin != "" {
		return c.PublicOrigin
	}
	if !validHost(r.Host) {
		return ""
	}
	return c.Scheme(r) + "://" + strings.ToLower(r.Host)
}

// Absolute turns a server-relative URL ("/download?...") into an absolute
// URL for r. Other URLs are returned unchanged, as is a relative URL when no
// origin can be determined.
func (c Config) Absolute(r *http.Request, u string) string {
	if !strings.HasPrefix(u, "/") || strings.HasPrefix(u, "//") {
		return u
	}
	origin := c.Origin(r)
	if origin == "" {
		return u
	}
	return origin + u
}

// validHost accepts host, host:port, [ipv6] and [ipv6]:port made of DNS name
// characters only.
func validHost(host string) bool {
	if host == "" || len(host) > 255 {
		return false
	}
	name := host
	if h, port, err := net.SplitHostPort(host); err == nil {
		if port == "" || strings.Trim(port, "0123456789") != "" {
			return false
		}
		name = h
	} else if strings.HasPrefix(host, "[") {
		if !strings.HasSuffix(host, "]") {
			return false
		}
		name = host[1 : len(host)-1]
	}
	if addr, err := netip.ParseAddr(name); err == nil {
		return addr.Zone() == ""
	}
	if strings.ContainsAny(name, "[]:") {
		return false
	}
	for _, label := range strings.Split(name, ".") {
		if label == "" || len(label) > 63 {
			return false
		}
		for _, ch := range label {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-' || ch == '_') {
				return false
			}
		}
	}
	return true
}

// AllowsDashboardOrigin reports whether origin may call /session
// cross-origin.
func (c Config) AllowsDashboardOrigin(origin string) bool {
	origin = strings.ToLower(origin)
	for _, allowed := range c.DashboardOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}
