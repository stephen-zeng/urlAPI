package data

import "fmt"

var (
	apiOpenai   = []string{"openai"}
	apiAlibaba  = []string{"alibaba"}
	apiDeekseek = []string{"deepseek"}
	apiOtherapi = []string{"otherapi"}
	security    = []string{"dashpwd", "blocklist", "allow"}
	txt         = []string{"txt", "txtrandomenabled", "txtsummaryenabled"}
	img         = []string{"img"}
	web         = []string{"web"}
)

func InitSetting() error {
	fmt.Println("init setting ...")
	return nil
}
