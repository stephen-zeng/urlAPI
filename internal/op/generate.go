package op

import (
	"encoding/json"
	"errors"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func GenerateTextImage(task model.Task, host string, skipDB bool) (model.Task, GenerateResult, error) {
	settings := database.SettingsStore.Get()
	if _, ok := database.PromptMap[task.Target]; ok {
		task.Target = settings.Prompts.Templates[task.Target]
	}
	if task.API == "" {
		task.API = settings.Text.GenerationAPI
	}
	provider, ok := settings.Providers.ByName(task.API)
	if !ok {
		return failTask(task, "Unknown API"), GenerateResult{URL: settings.Text.FallbackImageURL}, errors.New("text generation unknown API")
	}
	if task.Model == "" {
		task.Model = provider.TextModel
	}
	filter := TaskQueueFilter{Type: "txt.gen", Target: task.Target, API: task.API}
	return ExecuteCachedTask(task, filter, skipDB, func(task *model.Task) (GenerateResult, error) {
		return generateText(task, host, provider, settings.Prompts.GenerationContext)
	})
}

func GenerateImage(task model.Task, host string, skipDB bool) (model.Task, GenerateResult, error) {
	settings := database.SettingsStore.Get()
	if task.API == "" {
		task.API = settings.Image.API
	}
	provider, ok := settings.Providers.ByName(task.API)
	if !ok {
		return failTask(task, "Imggen Process invalid API"), GenerateResult{URL: settings.Image.FallbackImageURL}, errors.New("image generation invalid API")
	}
	if task.Model == "" {
		task.Model = provider.ImageModel
	}
	if task.Size == "" {
		task.Size = provider.ImageSize
	}
	filter := TaskQueueFilter{Type: "img.gen", Size: task.Size, Target: task.Target, API: task.API}
	return ExecuteCachedTask(task, filter, skipDB, func(task *model.Task) (GenerateResult, error) {
		return generateImage(task, host, provider)
	})
}

func GenerateWebImage(task model.Task, host string, skipDB bool) (model.Task, GenerateResult, error) {
	settings := database.SettingsStore.Get()
	task.API = util.WebTargetHost(task.Target)
	if task.API == "" {
		err := errors.New("web image invalid URL")
		return failTask(task, "Invalid URL"), GenerateResult{URL: settings.Web.FallbackImageURL}, err
	}
	filter := TaskQueueFilter{Type: "web.img", Target: task.Target, API: task.API}
	return ExecuteCachedTask(task, filter, skipDB, func(task *model.Task) (GenerateResult, error) {
		return generateWebImage(task, host)
	})
}

func GenerateRandom(task model.Task, skipDB bool) (model.Task, GenerateResult, error) {
	settings := database.SettingsStore.Get()
	if task.API == "" {
		task.API = settings.Random.DefaultAPI
	}
	result, err := generateRandom(&task)
	SaveTask(task, skipDB)
	if err != nil {
		return task, GenerateResult{URL: settings.Random.FallbackImageURL}, err
	}
	return task, result, nil
}

func DownloadImage(target string) ([]byte, string, error) {
	return downloadImage(target)
}

func failTask(task model.Task, msg string) model.Task {
	task.Status = "failed"
	task.Return = msg
	return task
}

func setTaskResult(task *model.Task, result GenerateResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return err
	}
	task.Return = string(body)
	task.Status = "success"
	return nil
}
