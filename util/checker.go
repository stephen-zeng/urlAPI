package util

import (
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
