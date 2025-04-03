package util

import (
	"regexp"
	"strings"
)

// 通配符检查
func WildcardChecker(strs *[]string, str *string) bool {
	for _, r := range *strs {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(r), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, *str)
		if err != nil {
			continue
		}
		if match || r == *str {
			return true
		}
	}
	return false
}

// 完全检查
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
