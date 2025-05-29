package handler

var task2settingName = map[string]string{
	"txt.gen": "txt",
	"img.gen": "img",
	"web.img": "web",
}

var expiredSettingPosition = map[string]int{
	"txt": 3,
	"web": 3,
	"img": 2,
}

type apiQuery struct {
	Prompt string `form:"prompt"`
	Model  string `form:"model"`
	API    string `form:"api"`
	More   string `form:"more"`
	User   string `form:"user"`
	Repo   string `form:"repo"`
	Img    string `form:"img"`
	URL    string `form:"url"`
	Size   string `form:"size"`
}

type txt struct{}
type img struct{}
type web struct{}
type rand struct{}
