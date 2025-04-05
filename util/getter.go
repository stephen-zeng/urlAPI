package util

import (
	"encoding/json"
	"github.com/mozillazg/go-pinyin"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func GetKeyList(oriList *[]string) [][]string {
	ret := make([][]string, 26)
	args := pinyin.NewArgs()
	for _, str := range *oriList {
		switch {
		case str[0] >= 'a' && str[0] <= 'z':
			ret[str[0]-'a'] = append(ret[str[0]-'a'], str)
		case str[0] >= 'A' && str[0] <= 'Z':
			ret[str[0]] = append(ret[str[0]-'A'], str)
		default:
			py := pinyin.Pinyin(str, args)
			if len(py) == 0 {
				continue
			}
			index := py[0][0][0] - 'a'
			ret[index] = append(ret[index], str)
		}
	}
	return ret
}

func GetKeyIndex(str string) int {
	args := pinyin.NewArgs()
	switch {
	case str[0] >= 'a' && str[0] <= 'z':
		return int(str[0] - 'a')
	case str[0] >= 'A' && str[0] <= 'Z':
		return int(str[0] - 'A')
	default:
		py := pinyin.Pinyin(str, args)
		if len(py) == 0 {
			return 0
		}
		return int(py[0][0][0] - 'a')
	}
}
