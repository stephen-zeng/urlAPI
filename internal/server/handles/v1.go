package handles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"urlAPI/internal/database"
	"urlAPI/internal/llm"
	"urlAPI/internal/model"
	"urlAPI/internal/server/middleware"
	"urlAPI/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responsesRequest struct {
	Model           string          `json:"model"`
	Input           json.RawMessage `json:"input"`
	Instructions    string          `json:"instructions,omitempty"`
	Stream          bool            `json:"stream,omitempty"`
	Temperature     float64         `json:"temperature,omitempty"`
	TopP            float64         `json:"top_p,omitempty"`
	MaxOutputTokens int             `json:"max_output_tokens,omitempty"`
	Tools           []json.RawMessage `json:"tools,omitempty"`
	ToolChoice      any             `json:"tool_choice,omitempty"`
	ResponseFormat  any             `json:"response_format,omitempty"`
}

type responseInputItem struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type responseToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

type responseTool struct {
	Type        string                `json:"type"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Parameters  any                   `json:"parameters,omitempty"`
	Function    *responseToolFunction `json:"function,omitempty"`
}

// ChatCompletionHandler handles /v1/chat/completions requests.
func ChatCompletionHandler(c *gin.Context) {
	var req llm.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChatCompletion] Invalid request body from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if len(req.Messages) == 0 {
		log.Printf("[ChatCompletion] Messages are required but not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "messages are required"})
		return
	}

	providerName := c.GetHeader("X-Provider")
	if providerName == "" {
		providerName = "otherapi"
	}
	log.Printf("[ChatCompletion] Request from %s | Provider: %s | Model: %s | Stream: %v | Messages: %d",
		c.ClientIP(), providerName, req.Model, req.Stream, len(req.Messages))

	providerConfig, ok := database.SettingsStore.Get().Providers.ByName(providerName)
	if !ok {
		log.Printf("[ChatCompletion] Unknown provider '%s' from %s", providerName, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	if providerConfig.Enabled == false {
		log.Printf("[ChatCompletion] Provider '%s' is disabled", providerName)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "provider disabled"})
		return
	}

	if req.Model == "" && providerConfig.TextModel != "" {
		req.Model = providerConfig.TextModel
		log.Printf("[ChatCompletion] Using default model '%s' for provider '%s'", req.Model, providerName)
	}
	if req.Model == "" {
		log.Printf("[ChatCompletion] Model is required but not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "model is required"})
		return
	}

	client, err := llm.NewProvider(providerConfig)
	if err != nil {
		log.Printf("[ChatCompletion] Failed to create provider client for '%s': %v", providerName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ChatCompletion] Provider client created for '%s' | Endpoint: %s", providerName, providerConfig.Endpoint)

	task := model.Task{
		UUID:   uuid.New().String(),
		Time:   time.Now(),
		IP:     c.ClientIP(),
		Type:   "txt.gen",
		Target: req.Messages[0].Content,
		API:    providerName,
		Model:  req.Model,
		Status: "running",
	}

	// Record API Key usage if authenticated
	if key, ok := middleware.GetAPIKey(c); ok {
		database.APIKeyStore.IncrementUsage(key.KeyHash)
	}

	ctx := context.Background()

	if req.Stream {
		if !client.IsStreamingSupported() {
			log.Printf("[ChatCompletion] Streaming not supported by provider '%s'", providerName)
			c.JSON(http.StatusBadRequest, gin.H{"error": "streaming not supported by provider"})
			return
		}

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		log.Printf("[ChatCompletion] Starting stream request to provider '%s'", providerName)
		stream, err := client.ChatCompletionStream(ctx, req)
		if err != nil {
			log.Printf("[ChatCompletion] Stream init failed for provider '%s': %v", providerName, err)
			task.Status = "failed"
			task.Return = err.Error()
			database.CreateTask(&task)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Stream(func(w io.Writer) bool {
			eventCount := 0
			for event := range stream {
				eventCount++
				if event.Error != nil {
					log.Printf("[ChatCompletion] Stream error from provider '%s': %v", providerName, event.Error)
					fmt.Fprintf(w, "data: %s\n\n", llm.ToJSON(llm.ChatCompletionResponse{
						Object: "chat.completion.chunk",
						Error:  &llm.Error{Message: event.Error.Error(), Type: "api_error"},
					}))
					return false
				}
				if event.Done {
					log.Printf("[ChatCompletion] Stream completed for provider '%s' | Events: %d", providerName, eventCount)
					fmt.Fprintf(w, "data: [DONE]\n\n")
					task.Status = "success"
					database.CreateTask(&task)
					return false
				}
				fmt.Fprintf(w, "data: %s\n\n", event.Data)
			}
			log.Printf("[ChatCompletion] Stream ended unexpectedly for provider '%s' | Events: %d", providerName, eventCount)
			return false
		})
		return
	}

	log.Printf("[ChatCompletion] Sending non-stream request to provider '%s'", providerName)
	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		log.Printf("[ChatCompletion] Request failed for provider '%s': %v", providerName, err)
		task.Status = "failed"
		task.Return = err.Error()
		database.CreateTask(&task)
		if resp != nil && resp.Error != nil {
			log.Printf("[ChatCompletion] Provider '%s' returned API error: %s", providerName, resp.Error.Message)
			c.JSON(http.StatusBadGateway, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ChatCompletion] Request succeeded for provider '%s' | Response ID: %s | Choices: %d | Tokens: %d",
		providerName, resp.ID, len(resp.Choices), resp.Usage.TotalTokens)
	task.Status = "success"
	body, _ := json.Marshal(resp)
	task.Return = string(body)
	database.CreateTask(&task)

	c.JSON(http.StatusOK, resp)
}

func ResponsesHandler(c *gin.Context) {
	var req responsesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Responses] Invalid request body from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	messages, err := responsesInputToMessages(req)
	if err != nil {
		log.Printf("[Responses] Invalid input from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tools := responsesToolsToChatTools(req.Tools)
	toolChoice := responsesToolChoice(req.ToolChoice, len(tools) > 0)

	chatReq := llm.ChatCompletionRequest{
		Model:          req.Model,
		Messages:       messages,
		Temperature:    req.Temperature,
		TopP:           req.TopP,
		MaxTokens:      req.MaxOutputTokens,
		Stream:         req.Stream,
		Tools:          tools,
		ToolChoice:     toolChoice,
		ResponseFormat: req.ResponseFormat,
	}

	providerName := c.GetHeader("X-Provider")
	if providerName == "" {
		providerName = "otherapi"
	}
	log.Printf("[Responses] Request from %s | Provider: %s | Model: %s | Stream: %v | Messages: %d",
		c.ClientIP(), providerName, chatReq.Model, chatReq.Stream, len(chatReq.Messages))

	providerConfig, ok := database.SettingsStore.Get().Providers.ByName(providerName)
	if !ok {
		log.Printf("[Responses] Unknown provider '%s' from %s", providerName, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}
	if !providerConfig.Enabled {
		log.Printf("[Responses] Provider '%s' is disabled", providerName)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "provider disabled"})
		return
	}
	if chatReq.Model == "" && providerConfig.TextModel != "" {
		chatReq.Model = providerConfig.TextModel
		log.Printf("[Responses] Using default model '%s' for provider '%s'", chatReq.Model, providerName)
	}
	if chatReq.Model == "" {
		log.Printf("[Responses] Model is required but not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "model is required"})
		return
	}

	client, err := llm.NewProvider(providerConfig)
	if err != nil {
		log.Printf("[Responses] Failed to create provider client for '%s': %v", providerName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if key, ok := middleware.GetAPIKey(c); ok {
		database.APIKeyStore.IncrementUsage(key.KeyHash)
	}

	ctx := context.Background()
	if chatReq.Stream {
		handleResponsesStream(c, client, chatReq, providerName)
		return
	}

	log.Printf("[Responses] Sending non-stream request to provider '%s'", providerName)
	resp, err := client.ChatCompletion(ctx, chatReq)
	if err != nil {
		log.Printf("[Responses] Request failed for provider '%s': %v", providerName, err)
		if resp != nil && resp.Error != nil {
			c.JSON(http.StatusBadGateway, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	outputText := chatCompletionText(resp)
	log.Printf("[Responses] Request succeeded for provider '%s' | Response ID: %s | Output length: %d",
		providerName, resp.ID, len(outputText))
	c.JSON(http.StatusOK, responsesBody(resp, outputText))
}

func handleResponsesStream(c *gin.Context, client llm.Provider, req llm.ChatCompletionRequest, providerName string) {
	if !client.IsStreamingSupported() {
		log.Printf("[Responses] Streaming not supported by provider '%s'", providerName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "streaming not supported by provider"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	log.Printf("[Responses] Starting stream request to provider '%s'", providerName)
	stream, err := client.ChatCompletionStream(context.Background(), req)
	if err != nil {
		log.Printf("[Responses] Stream init failed for provider '%s': %v", providerName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseID := "resp_" + uuid.New().String()
	itemID := "msg_" + uuid.New().String()
	eventCount := 0
	sequence := 0
	outputText := strings.Builder{}
	started := false
	c.Stream(func(w io.Writer) bool {
		if !started {
			writeResponsesOutputStart(w, responseID, itemID, req.Model, &sequence)
			started = true
			return true
		}

		event, ok := <-stream
		if !ok {
			log.Printf("[Responses] Stream closed without explicit done for provider '%s' | Events: %d", providerName, eventCount)
			writeResponsesOutputDone(w, itemID, outputText.String(), &sequence)
			fmt.Fprintf(w, "event: response.completed\ndata: %s\n\n", llm.ToJSON(responsesStreamCompletedBody(responseID, itemID, req.Model, outputText.String(), nextSequence(&sequence))))
			return false
		}
		if event.Error != nil {
			log.Printf("[Responses] Stream error from provider '%s': %v", providerName, event.Error)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", llm.ToJSON(gin.H{"error": event.Error.Error()}))
			return false
		}
		if event.Done {
			log.Printf("[Responses] Stream completed for provider '%s' | Events: %d", providerName, eventCount)
			writeResponsesOutputDone(w, itemID, outputText.String(), &sequence)
			fmt.Fprintf(w, "event: response.completed\ndata: %s\n\n", llm.ToJSON(responsesStreamCompletedBody(responseID, itemID, req.Model, outputText.String(), nextSequence(&sequence))))
			return false
		}
		delta := chatStreamDelta(event.Data)
		if delta == "" {
			return true
		}
		eventCount++
		outputText.WriteString(delta)
		fmt.Fprintf(w, "event: response.output_text.delta\ndata: %s\n\n", llm.ToJSON(gin.H{
			"type":          "response.output_text.delta",
			"sequence_number": nextSequence(&sequence),
			"item_id":       itemID,
			"output_index":  0,
			"content_index": 0,
			"delta":         delta,
			"logprobs":      []gin.H{},
		}))
		return true
	})
}

func responsesInputToMessages(req responsesRequest) ([]llm.ChatMessage, error) {
	messages := make([]llm.ChatMessage, 0)
	if strings.TrimSpace(req.Instructions) != "" {
		messages = append(messages, llm.ChatMessage{Role: "system", Content: req.Instructions})
	}

	var text string
	if err := json.Unmarshal(req.Input, &text); err == nil {
		text = strings.TrimSpace(text)
		if text == "" {
			return nil, fmt.Errorf("input is required")
		}
		return append(messages, llm.ChatMessage{Role: "user", Content: text}), nil
	}

	var items []responseInputItem
	if err := json.Unmarshal(req.Input, &items); err != nil {
		return nil, fmt.Errorf("input must be a string or message array")
	}
	for _, item := range items {
		role := item.Role
		if role == "" {
			role = "user"
		}
		content := responseContentText(item.Content)
		if strings.TrimSpace(content) == "" {
			continue
		}
		messages = append(messages, llm.ChatMessage{Role: role, Content: content})
	}
	if len(messages) == 0 || len(messages) == 1 && messages[0].Role == "system" {
		return nil, fmt.Errorf("input is required")
	}
	return messages, nil
}

func responseContentText(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		parts := make([]string, 0, len(v))
		for _, part := range v {
			m, ok := part.(map[string]any)
			if !ok {
				continue
			}
			if text, ok := m["text"].(string); ok {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n")
	default:
		return fmt.Sprintf("%v", content)
	}
}

func responsesToolsToChatTools(rawTools []json.RawMessage) []llm.Tool {
	if len(rawTools) == 0 {
		return nil
	}

	tools := make([]llm.Tool, 0, len(rawTools))
	for _, raw := range rawTools {
		var source responseTool
		if err := json.Unmarshal(raw, &source); err != nil {
			log.Printf("[Responses] Skipping malformed tool: %v", err)
			continue
		}

		toolType := source.Type
		name := source.Name
		description := source.Description
		parameters := source.Parameters
		if source.Function != nil {
			name = source.Function.Name
			description = source.Function.Description
			parameters = source.Function.Parameters
		}
		if toolType == "" && name != "" {
			toolType = "function"
		}

		if toolType != "function" {
			log.Printf("[Responses] Skipping unsupported tool type '%s'", toolType)
			continue
		}
		if strings.TrimSpace(name) == "" {
			log.Printf("[Responses] Skipping function tool with empty name")
			continue
		}

		var tool llm.Tool
		tool.Type = "function"
		tool.Function.Name = name
		tool.Function.Description = description
		tool.Function.Parameters = parameters
		tools = append(tools, tool)
	}

	if len(rawTools) > 0 {
		log.Printf("[Responses] Converted tools | Input: %d | Forwarded: %d", len(rawTools), len(tools))
	}
	return tools
}

func responsesToolChoice(choice any, hasTools bool) any {
	if choice == nil || !hasTools {
		return nil
	}
	value, ok := choice.(map[string]any)
	if !ok {
		return choice
	}
	if _, ok := value["function"]; ok {
		return choice
	}
	toolType, _ := value["type"].(string)
	name, _ := value["name"].(string)
	if toolType == "function" && strings.TrimSpace(name) != "" {
		return gin.H{"type": "function", "function": gin.H{"name": name}}
	}
	log.Printf("[Responses] Dropping unsupported tool_choice: %v", choice)
	return nil
}

func chatCompletionText(resp *llm.ChatCompletionResponse) string {
	if resp == nil || len(resp.Choices) == 0 {
		return ""
	}
	return resp.Choices[0].Message.Content
}

func responsesBody(resp *llm.ChatCompletionResponse, outputText string) gin.H {
	responseID := "resp_" + uuid.New().String()
	created := time.Now().Unix()
	modelName := ""
	usage := llm.Usage{}
	if resp != nil {
		if resp.ID != "" {
			responseID = resp.ID
		}
		if resp.Created != 0 {
			created = resp.Created
		}
		modelName = resp.Model
		usage = resp.Usage
	}
	return gin.H{
		"id":          responseID,
		"object":      "response",
		"created_at":  created,
		"status":      "completed",
		"model":       modelName,
		"output_text": outputText,
		"output": []gin.H{
			{
				"id":      "msg_" + uuid.New().String(),
				"type":    "message",
				"role":    "assistant",
				"status":  "completed",
				"content": []gin.H{{"type": "output_text", "text": outputText}},
			},
		},
		"usage": gin.H{
			"input_tokens":  usage.PromptTokens,
			"output_tokens": usage.CompletionTokens,
			"total_tokens":  usage.TotalTokens,
		},
	}
}

func nextSequence(sequence *int) int {
	value := *sequence
	*sequence = *sequence + 1
	return value
}

func responsesStreamCompletedBody(responseID string, itemID string, modelName string, outputText string, sequenceNumber int) gin.H {
	return gin.H{
		"type":            "response.completed",
		"sequence_number": sequenceNumber,
		"response": gin.H{
			"id":          responseID,
			"object":      "response",
			"created_at":  time.Now().Unix(),
			"status":      "completed",
			"model":       modelName,
			"output_text": outputText,
			"output": []gin.H{
				{
					"id":      itemID,
					"type":    "message",
					"role":    "assistant",
					"status":  "completed",
					"content": []gin.H{{"type": "output_text", "text": outputText, "annotations": []gin.H{}, "logprobs": []gin.H{}}},
				},
			},
			"usage": nil,
		},
	}
}

func writeResponsesOutputStart(w io.Writer, responseID string, itemID string, modelName string, sequence *int) {
	fmt.Fprintf(w, "event: response.created\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.created",
		"sequence_number": nextSequence(sequence),
		"response": gin.H{
			"id":         responseID,
			"object":     "response",
			"created_at": time.Now().Unix(),
			"status":     "in_progress",
			"model":      modelName,
		},
	}))
	fmt.Fprintf(w, "event: response.in_progress\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.in_progress",
		"sequence_number": nextSequence(sequence),
		"response": gin.H{
			"id":         responseID,
			"object":     "response",
			"created_at": time.Now().Unix(),
			"status":     "in_progress",
			"model":      modelName,
		},
	}))
	fmt.Fprintf(w, "event: response.output_item.added\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.output_item.added",
		"sequence_number": nextSequence(sequence),
		"output_index":    0,
		"item": gin.H{
			"id":      itemID,
			"type":    "message",
			"role":    "assistant",
			"status":  "in_progress",
			"content": []gin.H{},
		},
	}))
	fmt.Fprintf(w, "event: response.content_part.added\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.content_part.added",
		"sequence_number": nextSequence(sequence),
		"item_id":         itemID,
		"output_index":    0,
		"content_index":   0,
		"part":            gin.H{"type": "output_text", "text": "", "annotations": []gin.H{}, "logprobs": []gin.H{}},
	}))
}

func writeResponsesOutputDone(w io.Writer, itemID string, outputText string, sequence *int) {
	fmt.Fprintf(w, "event: response.output_text.done\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.output_text.done",
		"sequence_number": nextSequence(sequence),
		"item_id":         itemID,
		"output_index":    0,
		"content_index":   0,
		"text":            outputText,
		"logprobs":        []gin.H{},
	}))
	fmt.Fprintf(w, "event: response.content_part.done\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.content_part.done",
		"sequence_number": nextSequence(sequence),
		"item_id":         itemID,
		"output_index":    0,
		"content_index":   0,
		"part":            gin.H{"type": "output_text", "text": outputText, "annotations": []gin.H{}, "logprobs": []gin.H{}},
	}))
	fmt.Fprintf(w, "event: response.output_item.done\ndata: %s\n\n", llm.ToJSON(gin.H{
		"type":            "response.output_item.done",
		"sequence_number": nextSequence(sequence),
		"output_index":    0,
		"item": gin.H{
			"id":      itemID,
			"type":    "message",
			"role":    "assistant",
			"status":  "completed",
			"content": []gin.H{{"type": "output_text", "text": outputText, "annotations": []gin.H{}, "logprobs": []gin.H{}}},
		},
	}))
}

func chatStreamDelta(data string) string {
	var resp llm.ChatCompletionResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil || len(resp.Choices) == 0 {
		return ""
	}
	return resp.Choices[0].Delta.Content
}

// EmbeddingsHandler handles /v1/embeddings requests.
func EmbeddingsHandler(c *gin.Context) {
	var req llm.EmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Embeddings] Invalid request body from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	providerName := c.GetHeader("X-Provider")
	if providerName == "" {
		providerName = "otherapi"
	}
	log.Printf("[Embeddings] Request from %s | Provider: %s | Model: %s | Input length: %v",
		c.ClientIP(), providerName, req.Model, len(fmt.Sprintf("%v", req.Input)))

	providerConfig, ok := database.SettingsStore.Get().Providers.ByName(providerName)
	if !ok {
		log.Printf("[Embeddings] Unknown provider '%s' from %s", providerName, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	if providerConfig.Enabled == false {
		log.Printf("[Embeddings] Provider '%s' is disabled", providerName)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "provider disabled"})
		return
	}

	if req.Model == "" && providerConfig.EmbeddingModel != "" {
		req.Model = providerConfig.EmbeddingModel
		log.Printf("[Embeddings] Using default model '%s' for provider '%s'", req.Model, providerName)
	}

	client, err := llm.NewProvider(providerConfig)
	if err != nil {
		log.Printf("[Embeddings] Failed to create provider client for '%s': %v", providerName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !client.IsEmbeddingsSupported() {
		log.Printf("[Embeddings] Embeddings not supported by provider '%s'", providerName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "embeddings not supported by provider"})
		return
	}

	ctx := context.Background()
	log.Printf("[Embeddings] Sending request to provider '%s'", providerName)
	resp, err := client.Embeddings(ctx, req)
	if err != nil {
		log.Printf("[Embeddings] Request failed for provider '%s': %v", providerName, err)
		if resp != nil && resp.Error != nil {
			c.JSON(http.StatusBadGateway, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Embeddings] Request succeeded for provider '%s' | Data count: %d", providerName, len(resp.Data))
	c.JSON(http.StatusOK, resp)
}

// ModelsHandler handles /v1/models requests.
func ModelsHandler(c *gin.Context) {
	providerName := c.Query("provider")
	if providerName == "" {
		providerName = c.GetHeader("X-Provider")
	}
	log.Printf("[Models] Request from %s | Provider filter: %s", c.ClientIP(), providerName)

	var models []llm.Model

	if providerName != "" {
		providerConfig, ok := database.SettingsStore.Get().Providers.ByName(providerName)
		if !ok {
			log.Printf("[Models] Unknown provider '%s' from %s", providerName, c.ClientIP())
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
			return
		}

		client, err := llm.NewProvider(providerConfig)
		if err != nil {
			log.Printf("[Models] Failed to create provider client for '%s': %v", providerName, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		resp, err := client.Models(ctx)
		if err != nil {
			log.Printf("[Models] Failed to fetch models from provider '%s': %v", providerName, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		models = resp.Data
		log.Printf("[Models] Fetched %d models from provider '%s'", len(models), providerName)
	} else {
		settings := database.SettingsStore.Get()
		providers := map[string]util.ProviderConfig{
			"openai":    settings.Providers.OpenAI,
			"deepseek":  settings.Providers.DeepSeek,
			"alibaba":   settings.Providers.Alibaba,
			"anthropic": settings.Providers.Anthropic,
			"gemini":    settings.Providers.Gemini,
			"azure":     settings.Providers.Azure,
			"moonshot":  settings.Providers.Moonshot,
			"otherapi":  settings.Providers.OtherAPI,
		}

		for name, config := range providers {
			if !config.Enabled {
				continue
			}
			client, err := llm.NewProvider(config)
			if err != nil {
				log.Printf("[Models] Failed to create client for '%s': %v", name, err)
				continue
			}
			ctx := context.Background()
			resp, err := client.Models(ctx)
			if err != nil {
				log.Printf("[Models] Failed to fetch models from '%s': %v", name, err)
				continue
			}
			for _, m := range resp.Data {
				m.ID = name + "/" + m.ID
				models = append(models, m)
			}
		}
		log.Printf("[Models] Fetched total %d models from all providers", len(models))
	}

	c.JSON(http.StatusOK, llm.ModelsResponse{
		Object: "list",
		Data:   models,
	})
}
