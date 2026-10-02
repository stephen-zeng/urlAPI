package util

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func GetRegion(ip string) string {
	if value, ok := IPTmp[ip]; ok {
		return value
	}
	url := "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr?ip=" + ip
	resp, err := GlobalHTTPClient.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()
	jsonResp, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Unknown"
	}
	var response RegionResp
	err = json.Unmarshal(jsonResp, &response)
	if err != nil {
		return "Unknown"
	}

	var region string
	if response.Data.Country == "中国" {
		region = response.Data.Province
	} else {
		region = response.Data.Country
	}
	if len(IPTmp) >= 1000 {
		IPTmp = make(map[string]string)
	}

	IPTmp[ip] = region
	return region
}

func Downloader(url string) ([]byte, error) {
	resp, err := GlobalHTTPClient.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.WithMessage(err, resp.Status)
	}
	defer resp.Body.Close()
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	} else {
		return ret, nil
	}
}

func GetRepo(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var response []RepoContentResp
	if err = json.Unmarshal(jsonResponse, &response); err != nil {
		return nil, errors.WithStack(err)
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

func GetDate(ori string) time.Time {
	// yyyy.mm -> time.Time
	parts := strings.Split(ori, ".")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
}
