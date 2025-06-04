package util

import (
	"github.com/dlclark/regexp2"
	"image/png"
	"os"
	"regexp"
	"strings"
)

// 通配符检查，兼容正则表达式
func RegexChecker(strs *[]string, str *string) bool {
	for _, r := range *strs {
		if strings.HasPrefix(r, "re:") {
			pattern := r[3:]
			re := regexp2.MustCompile(pattern, 0)
			match, err := re.MatchString(*str)
			if err == nil && match {
				return true
			}
			continue
		}
		if strings.Contains(r, "*") {
			pattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(r), `\*`, ".*") + "$"
			re := regexp2.MustCompile(pattern, 0)
			match, err := re.MatchString(*str)
			if err == nil && match {
				return true
			}
			continue
		}
		if r == *str {
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

// 检查png图像是否合法，给路径
func PngChecker(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	if _, err = png.Decode(file); err != nil {
		return false
	}
	return true
}
