package img

type ImgResponse struct {
	URL          string `json:"url"`
	InitPrompt   string `json:"init_prompt"`
	ActualPrompt string `json:"actual_prompt"`
}
