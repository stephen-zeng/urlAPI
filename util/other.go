package util

import "strings"

func ListReplacer(list *[]string, old string, new string) {
	var ret []string
	for _, item := range *list {
		ret = append(ret, strings.Replace(item, old, new, -1))
	}
	*list = ret
}
