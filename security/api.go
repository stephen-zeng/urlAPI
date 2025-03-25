package security

import "urlAPI/util"

var (
	txt  = []string{"openai", "alibaba", "deepseek", "otherapi"}
	img  = []string{"openai", "alibaba"}
	rand = []string{"github", "gitee"}
	web  = []string{"arxiv", "bilibili", "ithome", "github", "gitee", "youtube"}
)

//func (*General) APIChecker() {}

func (info *TxtGen) APIChecker(general *General) {
	if !(util.ListChecker(&txt, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
	}
}

func (info *TxtSum) APIChecker(general *General) {
	if !(util.ListChecker(&txt, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
	}
}

func (info *ImgGen) APIChecker(general *General) {
	if !(util.ListChecker(&img, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
	}
}

func (info *WebImg) APIChecker(general *General) {
	if !(util.ListChecker(&web, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
	}
}

func (info *Rand) APIChecker(general *General) {
	if !(util.ListChecker(&rand, &(info.API))) {
		general.Info = "Invalid API"
		general.Unsafe = true
	}
}

func (info *General) APIChecker(general *General) {}
