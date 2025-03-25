package util

import (
	"encoding/json"
	"io"
	"net/url"
	"regexp"
	"strings"
)

// 只接受原始Referer
func RefererChecker(referers *[]string, referer *string) bool {
	domainParse, _ := url.Parse(*referer)
	domain := domainParse.Hostname()
	for _, r := range *referers {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(r), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, domain)
		if err != nil {
			continue
		}
		if match {
			return true
		}
	}
	return false
}

// api为空的时候返回正确
func ListChecker(apis *[]string, api *string) bool {
	if *api == "" {
		return true
	}
	for _, a := range *apis {
		if a == *api {
			return true
		}
	}
	return false
}

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
	if err != nil {
		return "Unknown"
	}
	var response RegionResp
	err = json.Unmarshal(jsonResp, &resp)
	if err != nil {
		return "Unknown"
	}
	return response.IPData.Info1
}
