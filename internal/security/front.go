package security

import "errors"

// Type: gen, sum, img, template, rand
// API: openai, alibaba, deepseek, otherapi, github, gitee
// Target: prompt, repositories
func NewRequest(data Config) error {
	var err error
	err = frequencyCheck(data.IP, data.Type, data.Target)
	if err != nil {
		return err
	}
	err = sourceCheck(data.Domain)
	if err != nil {
		return err
	}
	if data.Type == "txt.gen" || data.Type == "img.gen" {
		err = modelCheck(data.Type, data.API)
	}
	if err != nil {
		return err
	}
	switch data.Type {
	case "txt.gen":
		err = txtGenCheck(data.Target)
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
	return err
}
