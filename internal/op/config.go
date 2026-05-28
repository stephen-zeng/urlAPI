package op

import (
	"encoding/json"
	"os"
	"sync"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
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
	db        *database.SQLiteAdapter
	ImgPath   = "assets/img/"
	TaskQueue = SafeTaskQueue{
		Queue: make(map[TaskQueueFilter]TaskQueueItem),
	}
	TaskCounter = SafeTaskCounter{
		Counter: make(map[string]int),
	}
)

func Init() error {
	db = database.GetLocalDB()
	if err := os.RemoveAll(ImgPath); err != nil {
		return err
	}
	return os.MkdirAll(ImgPath, 0777)
}

func getEndpoint(api string) string {
	provider, ok := database.SettingsStore.Get().Providers.ByName(api)
	if !ok {
		return ""
	}
	return provider.Endpoint
}

type Session struct {
	// backend -> frontend
	SessionToken string          `json:"session_token"`
	SessionIP    string          `json:"session_ip"`
	SettingName  []string        `json:"setting_name"`
	SettingData  [][]string      `json:"setting_data"`
	SettingBody  json.RawMessage `json:"setting_body,omitempty"`
	TaskData     []model.Task    `json:"task_data"`
	TaskStats    TaskStats       `json:"task_stats"`
	TaskMaxPage  int             `json:"task_max_page"`
	RepoData     []model.Repo    `json:"repo_data"`

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

type GenerateResult struct {
	Prompt         string `json:"prompt"`
	OriginalPrompt string `json:"original_prompt"`
	ActualPrompt   string `json:"actual_prompt"`
	Response       string `json:"response"`
	URL            string `json:"url"`
}

type TaskStatItem struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type TaskStats map[string][]TaskStatItem

type TaskQueueItem struct {
	DB      model.Task `json:"db"`
	Return  GenerateResult
	Running bool `json:"running"`
}

type TaskQueueFilter struct {
	Type   string `json:"type"`
	Size   string `json:"size"`
	Target string `json:"target"`
	API    string `json:"api"`
}
