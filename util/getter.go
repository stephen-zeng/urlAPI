package util

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"regexp"
)

// 获取设备类型
func GetDeviceType(ua string) string {
	mobileRegexp := `(?i)(Mobile|Tablet|Android|iOS|iPhone|iPad|iPod)`
	desktopRegexp := `(?i)(Desktop|Windows|Macintosh|Linux)`
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
	url := "https://api.vore.top/api/IPdata?ip=" + ip
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
	return response.IPData.Info1
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
