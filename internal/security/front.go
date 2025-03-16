package security

import "errors"

// Type: gen, sum, img, template, rand
// API: openai, alibaba, deepseek, otherapi, github, gitee
// Target: prompt, repositories
func NewRequest(data Config) (string, error) {
	var err error
	err = frequencyCheck(data.IP, data.Type, data.Target)
	if err != nil {
		return "", err
	}
	info, err := sourceCheck(data.Domain)
	if err != nil {
		return "", err
	}
	if data.Type == "txt.gen" || data.Type == "img.gen" || data.Type == "txt.sum" {
		err = modelCheck(data.Type, data.API)
	}
	if err != nil {
		return "", err
	}
	switch data.Type {
	case "txt.gen":
		err = txtCheck(data.Target, data.Type)
	case "txt.sum":
		err = txtCheck(data.Target, data.Type)
	case "img.gen":
		err = imgGenCheck()
	case "rand":
		err = randCheck(data.API, data.Target)
	case "web.img":
		err = webImgCheck(data.API)
	case "download":
		err = nil
	default:
		err = errors.New("unknown type")
	}
	return info, err
}
