package processor

import (
	"os"
	"urlAPI/database"
)

var (
	ImgPath   = "/assets/img/"
	WebImgMap = map[string]string{
		"github.com":       "github",
		"gitee.com":        "gitee",
		"www.bilibili.com": "bilibili",
		"www.youtube.com":  "youtube",
		"arxiv.org":        "arxiv",
		"www.ithome.com":   "ithome",
	}
)

func init() {
	os.MkdirAll(ImgPath, 0777)
}

func getEndpoint(api string) string {
	switch api {
	case "openai":
		return database.SettingMap["openai"][5]
	case "alibaba":
		return "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	case "otherapi":
		return database.SettingMap["otherapi"][4]
	case "deepseek":
		return "https://api.deepseek.com/chat/completions"
	default:
		return ""
	}
}

type Interface interface {
	Process(data *database.Task) error
}

type Dashboard struct {
	// backend -> frontend
	SessionToken string          `json:"session_token"`
	SessionIP    string          `json:"session_ip"`
	SettingName  []string        `json:"setting_name"`
	SettingData  [][]string      `json:"setting_data"`
	TaskData     []database.Task `json:"task_data"`
	RepoData     []database.Repo `json:"repo_data"`

	// frontend -> backend
	Operation    string     `json:"operation"`
	LoginTerm    bool       `json:"login_term"`
	SettingEdit  [][]string `json:"setting_edit"`
	TaskCatagory string     `json:"task_catagory"`
	TaskBy       string     `json:"task_by"`
	RepoAPI      string     `json:"repo_api"`
	RepoInfo     string     `json:"repo_info"`
	RepoUUID     string     `json:"repo_uuid"`

	//both
	SettingPart string `json:"setting_part"`
}

type Download struct {
	Target      string `json:"target"`
	ReturnError string `json:"return_error"`
	Return      []byte `json:"return"`
}

type TxtGen struct {
	API    string `json:"api"`
	Model  string `json:"model"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Return string `json:"return"` // 这里是已经序列号好的json或者URL地址
	Host   string `json:"host"`
}

type TxtSum struct {
	TxtGen
}

type ImgGen struct {
	TxtGen
	Size string `json:"size"`
}

type WebImg struct {
	API    string `json:"api"`
	Target string `json:"target"`
	Host   string `json:"host"`
	Return string `json:"return"`
}

type Rand struct {
	API    string `json:"api"`
	Target string `json:"target"`
	Return string `json:"return"`
}
