package processor

import (
	"os"
	"sync"
	"urlAPI/database"
)

type SafeTaskQueue struct {
	Mu    sync.RWMutex
	Queue map[TaskQueueFilter]TaskQueueItem
}

type SafeTaskCounter struct {
	Mu      sync.RWMutex
	Counter map[string]int
}

var (
	ImgPath   = "assets/img/"
	TaskQueue = SafeTaskQueue{
		Queue: make(map[TaskQueueFilter]TaskQueueItem),
	}
	TaskCounter = SafeTaskCounter{
		Counter: make(map[string]int),
	}
)

func init() {
	os.RemoveAll(ImgPath)
	os.MkdirAll(ImgPath, 0777)

}

func getEndpoint(api string) string {
	switch api {
	case "openai":
		return database.SettingMap["openai"][5]
	case "alibaba":
		return database.SettingMap["alibaba"][5]
	case "otherapi":
		return database.SettingMap["otherapi"][3]
	case "deepseek":
		return database.SettingMap["deepseek"][3]
	default:
		return ""
	}
}

type Session struct {
	// backend -> frontend
	SessionToken string          `json:"session_token"`
	SessionIP    string          `json:"session_ip"`
	SettingName  []string        `json:"setting_name"`
	SettingData  [][]string      `json:"setting_data"`
	TaskData     []database.Task `json:"task_data"`
	TaskMaxPage  int             `json:"task_max_page"`
	RepoData     []database.Repo `json:"repo_data"`

	// frontend -> backend
	Operation    string     `json:"operation"`
	LoginTerm    bool       `json:"login_term"`
	SettingEdit  [][]string `json:"setting_edit"`
	TaskCatagory string     `json:"task_catagory"`
	TaskBy       string     `json:"task_by"`
	TaskPage     int        `json:"task_page"`
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
	Return string `json:"return"` // 这里是已经序列号好的json或者URL地址
	Host   string `json:"host"`
}

type ImgGen struct {
	API    string `json:"api"`
	Model  string `json:"model"`
	Target string `json:"target"`
	Return string `json:"return"` // 这里是已经序列号好的json或者URL地址
	Host   string `json:"host"`
	Size   string `json:"size"`
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

type TaskQueueItem struct {
	DB      database.Task `json:"db"`
	Return  string        `json:"return"`
	Running bool          `json:"running"`
}

type TaskQueueFilter struct {
	Type   string `json:"type"`
	Size   string `json:"size"`
	Target string `json:"target"`
	API    string `json:"api"`
}
