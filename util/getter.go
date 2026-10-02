package util

import (
	"context"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

// 获取设备类型
func GetDeviceType(ua string) string {
	mobileRegexp := `(?i)(Mobile|Tablet|Android|iOS|iPhone|iPad|iPod)`
	desktopRegexp := `(?i)(Desktop|Windows|Macintosh|Linux|PC)`
	botRegexp := `(?i)(Bot)`
	matched, _ := regexp.MatchString(mobileRegexp, ua)
	if matched {
		return "Mobile"
	}
	matched, _ = regexp.MatchString(desktopRegexp, ua)
	if matched {
		return "Desktop"
	}
	matched, _ = regexp.MatchString(botRegexp, ua)
	if matched {
		return "Bot"
	}
	return ""
}

// regionAPI resolves an IP address to a region; a variable for tests.
var regionAPI = "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr"

func GetRegion(ip string) string {
	if value, ok := ipRegions.get(ip); ok {
		return value
	}
	var response RegionResp
	err := doUpstreamJSON(context.Background(), upstreamRequest{
		URL:   regionAPI + "?" + url.Values{"ip": {ip}}.Encode(),
		Limit: 64 << 10,
	}, &response)
	if err != nil {
		return "Unknown"
	}

	var region string
	if response.Data.Country == "中国" {
		region = response.Data.Province
	} else {
		region = response.Data.Country
	}
	ipRegions.set(ip, region)
	return region
}

// Downloader fetches a URL and returns its body, failing on non-2xx
// responses and bodies larger than maxDownloadBytes.
func Downloader(url string) ([]byte, error) {
	return downloadContext(context.Background(), url)
}

func downloadContext(ctx context.Context, url string) ([]byte, error) {
	return doUpstream(ctx, upstreamRequest{URL: url, Limit: maxDownloadBytes})
}

func GetRepo(url string) ([]string, error) {
	var response []RepoContentResp
	if err := doUpstreamJSON(context.Background(), upstreamRequest{URL: url}, &response); err != nil {
		return nil, errors.WithMessage(err, "list repository")
	}
	var ret []string
	for _, repo := range response {
		ret = append(ret, repo.DownloadURL)
	}
	return ret, nil
}

func GetDomain(URL string) string {
	domainParse, err := url.Parse(URL)
	if err != nil {
		return ""
	}
	return domainParse.Hostname()
}
