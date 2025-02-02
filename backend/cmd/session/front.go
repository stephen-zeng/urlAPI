package session

import (
	"backend/internal/data"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

// string --> pwd, token
func Auth(dat Config) (string, error) {
	dash, err := data.FetchSetting(data.DataConfig(data.WithName([]string{"dash", "dashallowedip"})))
	if err != nil {
		return "", err
	}
	flag := false
	for _, item := range dash[1] {
		rgx := "^" + strings.ReplaceAll(regexp.QuoteMeta(item), `\*`, ".*") + "$"
		match, err := regexp.MatchString(rgx, dat.IP)
		if err != nil {
			continue
		}
		if match {
			flag = true
			break
		}
	}
	if !flag {
		log.Println("Authentication Failed by IP")
		return "", errors.New("Authentication Failed by IP")
	}
	if dat.Token == dash[0][0] {
		return "pwd", nil
	}
	tokens, err := data.FetchSession(data.DataConfig(data.WithToken(dat.Token)))
	if err != nil {
		log.Println(err)
		return "", errors.New("Token not found")
	}
	exp := tokens[0].Expire
	if time.Now().After(exp) {
		return "", errors.New("Token expired")
	} else {
		return "token", nil
	}
}
