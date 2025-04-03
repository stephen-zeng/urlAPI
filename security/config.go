package security

import "time"

type General struct {
	Referer string    `json:"referer"` //Complete Referer
	IP      string    `json:"ip"`
	Type    string    `json:"type"` // 任务类型
	Target  string    `json:"target"`
	Time    time.Time `json:"time"`
	Unsafe  bool      `json:"unsafe"`
	SkipDB  bool      `json:"skip_db"`
	Info    string    `json:"info"`
}

type TxtGen struct {
	API   string `json:"api"`
	Model string `json:"model"`
}

type TxtSum struct {
	API   string `json:"api"`
	Model string `json:"model"`
}

type ImgGen struct {
	API   string `json:"api"`
	Model string `json:"model"`
}

type Rand struct {
	API    string `json:"api"`
	Target string `json:"target"`
}

type WebImg struct {
	API string `json:"api"`
}

type FrequencyFilter struct {
	Type string `json:"type"`
	IP   string `json:"ip"`
}

type FrequencyData struct {
	Counter int       `json:"counter"`
	Time    time.Time `json:"time"`
}
