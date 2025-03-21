package security

type General struct {
	Referer string `json:"referer"` //Complete Referer
	IP      string `json:"ip"`
	Type    string `json:"type"`
	Target  string `json:"target"`
	Checked bool   `json:"checked"`
	SkipDB  bool   `json:"skip_db"`
}
type TxtGen struct {
	API    string `json:"api"`
	Model  string `json:"model"`
	Target string `json:"target"`
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
