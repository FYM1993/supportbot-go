package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// DashScopeClient 通义千问客户端
type DashScopeClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewDashScopeClient 创建通义千问客户端
func NewDashScopeClient(apiKey, model string, logger *zap.Logger) *DashScopeClient {
	return &DashScopeClient{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// Message 消息
type Message struct {
	Role       string      `json:"role"`                 // system, user, assistant, tool
	Content    string      `json:"content,omitempty"`    // 文本内容
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"` // 工具调用（assistant 角色）
	ToolCallID string      `json:"tool_call_id,omitempty"` // 工具调用ID（tool 角色）
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"` // "function"
	Function Function `json:"function"`
}

// Function 函数定义
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON 字符串
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model      string    `json:"model"`
	Input      Input     `json:"input"`
	Parameters Parameters `json:"parameters,omitempty"`
}

// Input 输入
type Input struct {
	Messages []Message `json:"messages"`
}

// Parameters 参数
type Parameters struct {
	Temperature float64                  `json:"temperature,omitempty"`
	TopP        float64                  `json:"top_p,omitempty"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Tools       []map[string]interface{} `json:"tools,omitempty"` // Function Calling 工具定义
	ToolChoice  string                   `json:"tool_choice,omitempty"` // "auto", "none"
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Output struct {
		Text         string     `json:"text"`
		FinishReason string     `json:"finish_reason"`
		Choices      []Choice   `json:"choices,omitempty"` // Function Calling 模式下的响应
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}

// Choice 选择项（Function Calling 模式）
type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

// Chat 调用通义千问聊天接口
func (c *DashScopeClient) Chat(messages []Message) (string, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	reqBody := ChatRequest{
		Model: c.model,
		Input: Input{
			Messages: messages,
		},
		Parameters: Parameters{
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API 返回错误: %d, body: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return chatResp.Output.Text, nil
}

// SimpleChat 简单聊天（单轮对话）
func (c *DashScopeClient) SimpleChat(systemPrompt, userMessage string) (string, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}
	return c.Chat(messages)
}

// ChatWithTools 支持工具调用的聊天
func (c *DashScopeClient) ChatWithTools(messages []Message, tools []map[string]interface{}) (*ChatResponse, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	reqBody := ChatRequest{
		Model: c.model,
		Input: Input{
			Messages: messages,
		},
		Parameters: Parameters{
			Temperature: 0.7,
			MaxTokens:   2000,
			Tools:       tools,
			ToolChoice:  "auto", // 自动决定是否调用工具
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	c.logger.Debug("ChatWithTools 请求", zap.String("request", string(jsonData)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	c.logger.Debug("ChatWithTools 响应", zap.String("response", string(body)))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API 返回错误: %d, body: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResp, nil
}

