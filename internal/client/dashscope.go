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
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
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
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Output struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
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

