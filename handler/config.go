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

type txt struct{}
type img struct{}
type web struct{}
type rand struct{}
