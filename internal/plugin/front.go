package plugin

import (
	"backend/internal/data"
	"errors"
)

var genContext, sumContext string

func Request(dat Config) (PluginResponse, error) {
	config, err := data.FetchSetting(data.DataConfig(data.WithSettingName([]string{"context"})))
	if err != nil {
		return PluginResponse{}, err
	}
	genContext = config[0][0]
	sumContext = config[0][1]
	if dat.API == "openai" {
		return openai(dat)
	} else if dat.API == "alibaba" {
		return alibaba(dat)
	} else if dat.API == "deepseek" {
		return deepseek(dat)
	} else if dat.API == "otherapi" {
		return otherapi(dat)
	} else if dat.API == "github" || dat.API == "gitee" {
		return random(dat.API, dat.Repo)
	} else {
		return PluginResponse{}, errors.New("No Valid API Option")
	}
}

func openai(dat Config) (PluginResponse, error) {
	model := dat.Model
	if dat.ImgPrompt != "" {
		prompt := dat.ImgPrompt
		size := dat.Size
		response, err := openaiImg(prompt, model, size)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else if dat.GenPrompt != "" {
		prompt := dat.GenPrompt
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

func alibaba(dat Config) (PluginResponse, error) {
	model := dat.Model
	if dat.ImgPrompt != "" {
		prompt := dat.ImgPrompt
		size := dat.Size
		response, err := alibabaImg(prompt, model, size)
		if err != nil {
			return PluginResponse{}, err
		} else {
			return response, nil
		}
	} else if dat.GenPrompt != "" {
		prompt := dat.GenPrompt
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

func deepseek(dat Config) (PluginResponse, error) {
	model := dat.Model
	if dat.GenPrompt != "" {
		prompt := dat.GenPrompt
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

func otherapi(dat Config) (PluginResponse, error) {
	model := dat.Model
	if dat.GenPrompt != "" {
		prompt := dat.GenPrompt
		contxt := genContext
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
