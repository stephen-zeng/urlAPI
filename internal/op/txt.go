package op

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func generateText(task *model.Task, provider util.ProviderConfig, context string) (GenerateResult, error) {
	endpoint := provider.Endpoint
	if endpoint == "" {
		task.Status = "failed"
		task.Return = "Unknown API"
		return GenerateResult{}, errors.New("text generation unknown API")
	}
	response, err := util.Txt(endpoint, provider.APIKey, task.Model, context, task.Target)
	if err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("text generation: %w", err)
	}
	img, err := util.DrawTxt(response)
	if err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("text image rendering: %w", err)
	}
	file, err := os.Create(ImgPath + task.UUID + ".png")
	if err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("text image create: %w", err)
	}
	defer file.Close()
	if _, err = io.Copy(file, bytes.NewReader(img)); err != nil {
		task.Status = "failed"
		task.Return = err.Error()
		return GenerateResult{}, fmt.Errorf("text image write: %w", err)
	}
	result := GenerateResult{Prompt: task.Target, Response: response, URL: downloadURL(task.UUID)}
	if err := setTaskResult(task, result); err != nil {
		return GenerateResult{}, err
	}
	return result, nil
}
