package processor

import "urlAPI/database"

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

type API struct {
	IP      string `json:"ip"`
	Target  string `json:"target"`
	API     string `json:"api"`
	Referer string `json:"referer"`
	Device  string `json:"device"`
	Size    string `json:"size"`
	Type    string `json:"type"`
}
