package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"urlAPI/util"
)

// openAICompatibleProvider implements Provider for OpenAI-compatible APIs.
type openAICompatibleProvider struct {
	config util.ProviderConfig
	name   string
}

func newOpenAICompatibleProvider(config util.ProviderConfig, name string) *openAICompatibleProvider {
	return &openAICompatibleProvider{
		config: config,
		name:   name,
	}
}

func (p *openAICompatibleProvider) Name() string {
	return p.name
}

func (p *openAICompatibleProvider) IsStreamingSupported() bool {
	return true
}

func (p *openAICompatibleProvider) IsEmbeddingsSupported() bool {
	return p.config.EmbeddingModel != ""
}

func (p *openAICompatibleProvider) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if p.config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint not configured")
	}
	if p.config.APIKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[OpenAIProvider] Sending request to %s | Model: %s | Stream: %v | Messages: %d",
		p.config.Endpoint, req.Model, req.Stream, len(req.Messages))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("User-Agent", "urlAPI/1.0")

	for k, v := range p.config.CustomHeaders {
		httpReq.Header.Set(k, v)
	}

	resp, err := util.GlobalHTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("[OpenAIProvider] HTTP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[OpenAIProvider] HTTP error %d from %s | Body: %s", resp.StatusCode, p.config.Endpoint, string(body))
		var errResp struct {
			Error Error `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
			return &ChatCompletionResponse{Error: &errResp.Error}, fmt.Errorf("API error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[OpenAIProvider] Failed to parse response: %v | Body: %s", err, string(body))
		return nil, err
	}

	log.Printf("[OpenAIProvider] Response received | Choices: %d | Usage: %d tokens",
		len(result.Choices), result.Usage.TotalTokens)

	return &result, nil
}

func (p *openAICompatibleProvider) ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<-chan StreamEvent, error) {
	if p.config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint not configured")
	}
	if p.config.APIKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}

	req.Stream = true
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[OpenAIProvider] Starting stream request to %s | Model: %s", p.config.Endpoint, req.Model)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.Header.Set("Connection", "keep-alive")
	httpReq.Header.Set("User-Agent", "urlAPI/1.0")

	for k, v := range p.config.CustomHeaders {
		httpReq.Header.Set(k, v)
	}

	resp, err := util.GlobalHTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("[OpenAIProvider] Stream HTTP request failed: %v", err)
		return nil, err
	}

	events := make(chan StreamEvent)
	go func() {
		defer close(events)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("[OpenAIProvider] Stream HTTP error %d from %s | Body: %s", resp.StatusCode, p.config.Endpoint, string(body))
			events <- StreamEvent{Error: fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))}
			return
		}

		log.Printf("[OpenAIProvider] Stream connection established to %s", p.config.Endpoint)
		scanner := bufio.NewScanner(resp.Body)
		eventCount := 0
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				log.Printf("[OpenAIProvider] Stream completed | Events: %d", eventCount)
				events <- StreamEvent{Done: true}
				return
			}
			eventCount++
			events <- StreamEvent{Data: data}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("[OpenAIProvider] Stream scanner error: %v | Events: %d", err, eventCount)
			events <- StreamEvent{Error: err}
		}
		log.Printf("[OpenAIProvider] Stream ended | Events: %d", eventCount)
	}()

	return events, nil
}

func (p *openAICompatibleProvider) Embeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	if p.config.EmbeddingModel == "" {
		return nil, fmt.Errorf("embedding model not configured")
	}

	embeddingEndpoint := strings.Replace(p.config.Endpoint, "/chat/completions", "/embeddings", 1)
	if embeddingEndpoint == p.config.Endpoint {
		embeddingEndpoint = p.config.Endpoint + "/embeddings"
	}

	req.Model = p.config.EmbeddingModel
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[OpenAIProvider] Sending embeddings request to %s | Model: %s", embeddingEndpoint, req.Model)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", embeddingEndpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("User-Agent", "urlAPI/1.0")

	for k, v := range p.config.CustomHeaders {
		httpReq.Header.Set(k, v)
	}

	resp, err := util.GlobalHTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("[OpenAIProvider] Embeddings HTTP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[OpenAIProvider] Embeddings HTTP error %d from %s | Body: %s", resp.StatusCode, embeddingEndpoint, string(body))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result EmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[OpenAIProvider] Failed to parse embeddings response: %v | Body: %s", err, string(body))
		return nil, err
	}

	log.Printf("[OpenAIProvider] Embeddings response received | Data count: %d", len(result.Data))

	return &result, nil
}

func (p *openAICompatibleProvider) Models(ctx context.Context) (*ModelsResponse, error) {
	modelsEndpoint := strings.Replace(p.config.Endpoint, "/chat/completions", "/models", 1)
	if modelsEndpoint == p.config.Endpoint {
		modelsEndpoint = p.config.Endpoint + "/models"
	}

	log.Printf("[OpenAIProvider] Fetching models from %s", modelsEndpoint)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", modelsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "urlAPI/1.0")

	for k, v := range p.config.CustomHeaders {
		httpReq.Header.Set(k, v)
	}

	resp, err := util.GlobalHTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("[OpenAIProvider] Models HTTP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[OpenAIProvider] Models HTTP error %d from %s | Body: %s", resp.StatusCode, modelsEndpoint, string(body))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result ModelsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[OpenAIProvider] Failed to parse models response: %v | Body: %s", err, string(body))
		return nil, err
	}

	log.Printf("[OpenAIProvider] Models response received | Count: %d", len(result.Data))

	return &result, nil
}
