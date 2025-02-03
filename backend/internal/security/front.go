package security

// Type: gen, sum, img, web, rand
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
		if err != nil {
			return err
		}
	}
	if data.Type == "rand" {
		err = repoCheck(data.API, data.User, data.Repo)
		if err != nil {
			return err
		}
	}
	if data.Type == "txt.gen" {
		err = txtGenCheck(data.Target)
		if err != nil {
			return err
		}
	}
	if data.Type == "img.gen" {
		err = imgGenCheck()
		if err != nil {
			return err
		}
	}
	return nil
}
