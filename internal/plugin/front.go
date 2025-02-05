package plugin

import (
	"errors"
)

const (
	genContext string = "You are a helpful assistant and need to give some sentence based on the prompt. Do not greet and give the answer directly."
	//	sumContext string = "You are a helpful assistant and need to summarize the text from the prompt. Do not greet and give the answer directly."
)

func Request(data Config) (PluginResponse, error) {
	if data.API == "openai" {
		return openai(data)
	} else if data.API == "alibaba" {
		return alibaba(data)
	} else if data.API == "deepseek" {
		return deepseek(data)
	} else if data.API == "otherapi" {
		return otherapi(data)
	} else if data.API == "github" || data.API == "gitee" {
		return random(data.API, data.Repo)
	} else {
		return PluginResponse{}, errors.New("No Valid API Option")
	}
}

func openai(data Config) (PluginResponse, error) {
	model := data.Model
	if data.ImgPrompt != "" {
		prompt := data.ImgPrompt
		size := data.Size
		response, err := openaiImg(prompt, model, size)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := genContext
		response, err := openaiTxt(prompt, contxt, model)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else {
		return PluginResponse{}, errors.New("No Valid Prompt")
	}
}

func alibaba(data Config) (PluginResponse, error) {
	model := data.Model
	if data.ImgPrompt != "" {
		prompt := data.ImgPrompt
		size := data.Size
		response, err := alibabaImg(prompt, model, size)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := genContext
		response, err := alibabaTxt(prompt, contxt, model)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else {
		return PluginResponse{}, errors.New("No Valid Prompt")
	}
}

func deepseek(data Config) (PluginResponse, error) {
	model := data.Model
	if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := genContext
		response, err := deepseekTxt(prompt, contxt, model)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else {
		return PluginResponse{}, errors.New("No Valid Prompt")
	}
}

func otherapi(data Config) (PluginResponse, error) {
	model := data.Model
	if data.GenPrompt != "" {
		prompt := data.GenPrompt
		contxt := "You are a helpful assistant and need to give some sentence based on the prompt"
		response, err := otherapiTxt(prompt, contxt, model)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else {
		return PluginResponse{}, errors.New("No Valid Prompt")
	}
}
