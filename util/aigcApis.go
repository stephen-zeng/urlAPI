package util

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Alibaba DashScope endpoints and polling parameters; variables for tests.
var (
	alibabaImageSynthesisAPI = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis"
	alibabaTaskAPI           = "https://dashscope.aliyuncs.com/api/v1/tasks/"
	// alibabaPollTimeout bounds the whole asynchronous image job, including
	// submission, polling and downloading the result.
	alibabaPollTimeout  = 2 * time.Minute
	alibabaPollInterval = time.Second
)

func bearer(token string) map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}
}

func Txt(endpoint, token, model, systemPrompt, prompt string) (string, error) {
	if endpoint == "" || token == "" || model == "" || systemPrompt == "" || prompt == "" {
		return "", errors.WithStack(errors.New("Util TxtAPI insufficient info"))
	}
	userMessage := TxtMessage{
		Role:    "user",
		Content: prompt,
	}
	developerMessage := TxtMessage{
		Role:    "system",
		Content: systemPrompt,
	}
	txtPayload := TxtPayload{
		Model:    model,
		Messages: []TxtMessage{developerMessage, userMessage},
	}
	jsonPayload, err := json.Marshal(txtPayload)
	if err != nil {
		return "", errors.WithStack(err)
	}
	var txtResp TxtResp
	err = doUpstreamJSON(context.Background(), upstreamRequest{
		Method:  "POST",
		URL:     endpoint,
		Body:    jsonPayload,
		Headers: bearer(token),
		Secrets: []string{token},
	}, &txtResp)
	if err != nil {
		return "", errors.WithMessage(err, "text generation")
	}
	if len(txtResp.Choices) == 0 {
		return "", errors.New("text generation: upstream returned no choices")
	}
	return txtResp.Choices[0].Message.Content, nil
}

func AlibabaImg(token, prompt, model, size string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), alibabaPollTimeout)
	defer cancel()
	return alibabaImg(ctx, token, prompt, model, size)
}

func alibabaImg(ctx context.Context, token, prompt, model, size string) ([]byte, string, error) {
	imgPayload := AlibabaImgPayload{
		Model:      model,
		Input:      AlibabaImgInput{Prompt: prompt},
		Parameters: AlibabaImgParameters{Size: size, N: 1},
	}
	jsonPayload, err := json.Marshal(imgPayload)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	headers := bearer(token)
	headers["X-DashScope-Async"] = "enable"
	var response AlibabaImgResp
	err = doUpstreamJSON(ctx, upstreamRequest{
		Method:  "POST",
		URL:     alibabaImageSynthesisAPI,
		Body:    jsonPayload,
		Headers: headers,
		Secrets: []string{token},
	}, &response)
	if err != nil {
		return nil, "", errors.WithMessage(err, "alibaba image submit")
	}
	id := response.Output.TaskID
	if id == "" {
		return nil, "", errors.New("alibaba image submit: upstream returned no task id")
	}

	ticker := time.NewTicker(alibabaPollInterval)
	defer ticker.Stop()
	for status := response.Output.TaskStatus; status == "PENDING" || status == "RUNNING"; status = response.Output.TaskStatus {
		select {
		case <-ctx.Done():
			return nil, "", errors.Wrapf(ctx.Err(), "alibaba image task %s did not finish (last status %s)", id, status)
		case <-ticker.C:
		}
		response = AlibabaImgResp{}
		err = doUpstreamJSON(ctx, upstreamRequest{
			URL:     alibabaTaskAPI + url.PathEscape(id),
			Headers: map[string]string{"Authorization": "Bearer " + token},
			Secrets: []string{token},
		}, &response)
		if err != nil {
			return nil, "", errors.WithMessage(err, "alibaba image poll")
		}
	}

	if response.Output.TaskStatus != "SUCCEEDED" {
		return nil, "", errors.Errorf("alibaba image task %s ended with status %q", id, response.Output.TaskStatus)
	}
	if len(response.Output.Results) == 0 || response.Output.Results[0].URL == "" {
		return nil, "", errors.New("alibaba image task returned no results")
	}
	actualPrompt := response.Output.Results[0].ActualPrompt
	ret, err := downloadContext(ctx, response.Output.Results[0].URL)
	if err != nil {
		return nil, "", errors.WithMessage(err, "alibaba image download")
	}
	return ret, actualPrompt, nil
}

func OpenaiImg(endpoint, token, prompt, model, size string) ([]byte, error) {
	imgPayload := OpenaiImgPayload{
		Model:  model,
		Prompt: prompt,
		Size:   size,
		N:      1,
	}
	jsonPayload, err := json.Marshal(imgPayload)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var response OpenaiImgResp
	err = doUpstreamJSON(context.Background(), upstreamRequest{
		Method:  "POST",
		URL:     endpoint,
		Body:    jsonPayload,
		Headers: bearer(token),
		Secrets: []string{token},
	}, &response)
	if err != nil {
		return nil, errors.WithMessage(err, "openai image")
	}
	if len(response.Data) == 0 || response.Data[0].URL == "" {
		return nil, errors.New("openai image: upstream returned no images")
	}
	ret, err := Downloader(response.Data[0].URL)
	if err != nil {
		return nil, errors.WithMessage(err, "openai image download")
	}
	return ret, nil
}
