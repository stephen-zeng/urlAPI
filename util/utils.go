package utils

import (
	"regexp"
	"strings"
)

func DomainChecker(referers *[]string, domain *string) bool {
	for _, referer := range *referers {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(referer), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, *domain)
		if err != nil {
			continue
		}
		if match {
			return true
		}
	}
	return false
}
