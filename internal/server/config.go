package server

type sessionResp struct {
	Operation    string     `json:"operation"`
	LoginTerm    bool       `json:"login_term"`
	TaskBy       string     `json:"task_by"`
	TaskCatagory string     `json:"task_catagory"`
	SettingPart  string     `json:"setting_part"`
	SettingEdit  [][]string `json:"setting_edit"`
	RepoAPI      string     `json:"repo_api"`
	RepoInfo     string     `json:"repo_info"`
	RepoUUID     string     `json:"repo_uuid"`
}
