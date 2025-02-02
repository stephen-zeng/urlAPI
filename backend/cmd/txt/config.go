package txt

type TxtResponse struct {
	Response string `json:"response"`
	Context  string `json:"context"`
	Prompt   string `json:"prompt"`
}
