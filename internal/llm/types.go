package llm

import (
	"context"
	"encoding/json"
	"io"
)

/** @brief 聊天补全中的单条消息。 */
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

/** @brief 聊天消息中的工具调用。 */
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

/** @brief 聊天补全接口请求结构。 */
type ChatCompletionRequest struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	Temperature      float64       `json:"temperature,omitempty"`
	TopP             float64       `json:"top_p,omitempty"`
	MaxTokens        int           `json:"max_tokens,omitempty"`
	PresencePenalty  float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64       `json:"frequency_penalty,omitempty"`
	Stream           bool          `json:"stream,omitempty"`
	Stop             []string      `json:"stop,omitempty"`
	Seed             int           `json:"seed,omitempty"`
	Tools            []Tool        `json:"tools,omitempty"`
	ToolChoice       any           `json:"tool_choice,omitempty"`
	ResponseFormat   any           `json:"response_format,omitempty"`
}

/** @brief 可供模型调用的工具定义。 */
type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	} `json:"function"`
}

/** @brief 聊天补全接口响应结构。 */
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Error   *Error   `json:"error,omitempty"`
}

/** @brief 聊天补全响应中的单个候选结果。 */
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Delta        ChatMessage `json:"delta,omitempty"`
}

/** @brief Token 使用量统计。 */
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

/** @brief LLM 接口错误响应。 */
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

/** @brief Embeddings 接口请求结构。 */
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input any    `json:"input"`
	User  string `json:"user,omitempty"`
}

/** @brief Embeddings 接口响应结构。 */
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
	Error  *Error      `json:"error,omitempty"`
}

/** @brief 单条向量结果。 */
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

/** @brief 模型元数据信息。 */
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

/** @brief 模型列表接口响应结构。 */
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
	Error  *Error  `json:"error,omitempty"`
}

/** @brief SSE 流中的单个事件。 */
type StreamEvent struct {
	Data  string
	Error error
	Done  bool
}

/** @brief LLM 提供方统一接口。 */
type Provider interface {
	// Name returns the provider name.
	Name() string

	// ChatCompletion sends a chat completion request and returns the response.
	ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)

	// ChatCompletionStream sends a chat completion request and returns a stream of events.
	ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<-chan StreamEvent, error)

	// Embeddings sends an embedding request and returns the response.
	Embeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error)

	// Models returns the list of available models.
	Models(ctx context.Context) (*ModelsResponse, error)

	// IsStreamingSupported returns whether the provider supports streaming.
	IsStreamingSupported() bool

	// IsEmbeddingsSupported returns whether the provider supports embeddings.
	IsEmbeddingsSupported() bool
}

/** @brief SSE 事件写入辅助器。 */
type SSEWriter struct {
	Writer io.Writer
}

/**
 * @brief 写入单个 SSE 数据事件。
 * @param data 事件数据。
 * @return error 写入失败时返回错误。
 */
func (w *SSEWriter) WriteEvent(data string) error {
	_, err := w.Writer.Write([]byte("data: " + data + "\n\n"))
	return err
}

/**
 * @brief 写入 SSE 结束事件。
 * @return error 写入失败时返回错误。
 */
func (w *SSEWriter) WriteDone() error {
	_, err := w.Writer.Write([]byte("data: [DONE]\n\n"))
	return err
}

/**
 * @brief 将任意值序列化为 JSON 字符串。
 * @param v 待序列化的值。
 * @return string JSON 字符串。
 */
func ToJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
