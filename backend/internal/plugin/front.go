package plugin

import "errors"

func Require(data Config) (string, error) {
	if data.API == "openai" {
		return openai(data)
	} else if data.API == "alibaba" {
		return alibaba(data)
	} else if data.API == "deepseek" {
		return deepseek(data)
	} else if data.API == "otherapi" {
		return otherapi(data)
	} else {
		return "", errors.New("plugin.response.error")
	}
}

func openai(data Config) (string, error) {
	model := data.Model
	if data.ImgPrompt != "" {
		prompt := data.ImgPrompt
		size := data.Size
		response, err := openaiImg(prompt, model, size)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := "You are a helpful assistant and need to give some sentence based on the prompt"
		response, err := openaiTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else if data.SumPrompt != "" {
		prompt := data.SumPrompt
		contxt := "You are a helpful assistant and need to summarize the text from the prompt"
		response, err := openaiTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else {
		return "", errors.New("plugin.response.error")
	}
}

func alibaba(data Config) (string, error) {
	model := data.Model
	if data.ImgPrompt != "" {
		prompt := data.ImgPrompt
		size := data.Size
		response, err := alibabaImg(prompt, model, size)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := "You are a helpful assistant and need to give some sentence based on the prompt"
		response, err := alibabaTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else if data.SumPrompt != "" {
		prompt := data.SumPrompt
		contxt := "You are a helpful assistant and need to summarize the text from the prompt"
		response, err := alibabaTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else {
		return "", errors.New("plugin.response.error")
	}
}

func deepseek(data Config) (string, error) {
	model := data.Model
	if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := "You are a helpful assistant and need to give some sentence based on the prompt"
		response, err := deepseekTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else if data.SumPrompt != "" {
		prompt := data.SumPrompt
		contxt := "You are a helpful assistant and need to summarize the text from the prompt"
		response, err := deepseekTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else {
		return "", errors.New("plugin.response.error")
	}
}

func otherapi(data Config) (string, error) {
	model := data.Model
	if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := "You are a helpful assistant and need to give some sentence based on the prompt"
		response, err := otherapiTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else if data.SumPrompt != "" {
		prompt := data.SumPrompt
		contxt := "You are a helpful assistant and need to summarize the text from the prompt"
		response, err := otherapiTxt(prompt, contxt, model)
		if err != nil {
			return "", err
		} else {
			return response, nil
		}
	} else {
		return "", errors.New("plugin.response.error")
	}
}
