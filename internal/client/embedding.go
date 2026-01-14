package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// EmbeddingClient 通义千问 Embedding 客户端
type EmbeddingClient struct {
	apiKey string
	logger *zap.Logger
	client *http.Client
}

// EmbeddingRequest 请求结构
type EmbeddingRequest struct {
	Model      string                 `json:"model"`
	Input      EmbeddingInput         `json:"input"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// EmbeddingInput 输入结构（支持单个文本或文本数组）
type EmbeddingInput struct {
	Texts []string `json:"texts"`
}

// EmbeddingResponse 响应结构
type EmbeddingResponse struct {
	Output struct {
		Embeddings []struct {
			TextIndex int       `json:"text_index"`
			Embedding []float64 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}

// NewEmbeddingClient 创建 Embedding 客户端
func NewEmbeddingClient(apiKey string, logger *zap.Logger) *EmbeddingClient {
	return &EmbeddingClient{
		apiKey: apiKey,
		logger: logger,
		client: &http.Client{},
	}
}

// GetEmbedding 获取单个文本的向量
func (c *EmbeddingClient) GetEmbedding(text string) ([]float64, error) {
	embeddings, err := c.GetEmbeddings([]string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return embeddings[0], nil
}

// GetEmbeddings 批量获取文本向量
func (c *EmbeddingClient) GetEmbeddings(texts []string) ([][]float64, error) {
	c.logger.Info("获取文本向量", zap.Int("count", len(texts)))

	// 构建请求
	reqBody := EmbeddingRequest{
		Model: "text-embedding-v2", // 通义千问 Embedding 模型
		Input: EmbeddingInput{
			Texts: texts,
		},
		Parameters: map[string]interface{}{
			"text_type": "document", // document 或 query
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var embResp EmbeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取向量
	embeddings := make([][]float64, len(embResp.Output.Embeddings))
	for i, emb := range embResp.Output.Embeddings {
		embeddings[emb.TextIndex] = emb.Embedding
	}

	c.logger.Info("向量获取成功",
		zap.Int("count", len(embeddings)),
		zap.Int("dimension", len(embeddings[0])),
		zap.Int("tokens", embResp.Usage.TotalTokens))

	return embeddings, nil
}

// GetQueryEmbedding 获取查询文本的向量（与文档向量略有不同）
func (c *EmbeddingClient) GetQueryEmbedding(query string) ([]float64, error) {
	c.logger.Info("获取查询向量", zap.String("query", query))

	reqBody := EmbeddingRequest{
		Model: "text-embedding-v2",
		Input: EmbeddingInput{
			Texts: []string{query},
		},
		Parameters: map[string]interface{}{
			"text_type": "query", // 查询模式
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var embResp EmbeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(embResp.Output.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embResp.Output.Embeddings[0].Embedding, nil
}
