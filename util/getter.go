package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"math/big"
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
	if value, ok := IPTmp[ip]; ok {
		return value
	}
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
	if len(IPTmp) >= 1000 {
		IPTmp = make(map[string]string)
	}
	IPTmp[ip] = response.IPData.Info1
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

func GetRandomString() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	randomNumber := n.String()
	hash := sha256.Sum256([]byte(randomNumber))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func GetShortRandomString(len int) string {
	if len >= 64 {
		return GetRandomString()
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	randomNumber := n.String()
	hash := sha256.Sum256([]byte(randomNumber))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr[:len]
}
