package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/supportbot/supportbot-go/internal/client"
	"github.com/supportbot/supportbot-go/internal/model"
	"go.uber.org/zap"
)

// ClassifierService 问题分类服务
type ClassifierService struct {
	llmClient   *client.DashScopeClient
	redisClient *redis.Client
	httpClient  *http.Client
	agentURLs   map[string]string // category -> agent URL
	logger      *zap.Logger
	systemPrompt string
	categories   map[string]CategoryInfo
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	AgentURL    string   `json:"agentUrl"`
}

// NewClassifierService 创建问题分类服务
func NewClassifierService(
	llmClient *client.DashScopeClient,
	redisClient *redis.Client,
	categories map[string]CategoryInfo,
	systemPrompt string,
	logger *zap.Logger,
) *ClassifierService {
	agentURLs := make(map[string]string)
	for category, info := range categories {
		agentURLs[category] = info.AgentURL
	}

	return &ClassifierService{
		llmClient:    llmClient,
		redisClient:  redisClient,
		httpClient:   &http.Client{},
		agentURLs:    agentURLs,
		categories:   categories,
		systemPrompt: systemPrompt,
		logger:       logger,
	}
}

// ClassifyAndRoute 分类并路由问题
func (s *ClassifierService) ClassifyAndRoute(userID int64, question string) (*model.ClassifyResponse, error) {
	s.logger.Info("开始分类问题",
		zap.Int64("userId", userID),
		zap.String("question", question))

	// 1. 从 Redis 获取对话历史
	ctx := context.Background()
	historyKey := fmt.Sprintf("chat_history:%d", userID)
	history, _ := s.redisClient.LRange(ctx, historyKey, -5, -1).Result()

	// 2. 构建分类提示词
	prompt := s.buildClassifyPrompt(question, history)

	// 3. 调用 LLM 分类
	response, err := s.llmClient.SimpleChat(s.systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM 分类失败: %w", err)
	}

	// 4. 解析分类结果
	category := s.parseCategory(response)
	
	result := &model.ClassifyResponse{
		Category:    category,
		Confidence:  0.9,
		Description: s.categories[category].Description,
	}

	s.logger.Info("问题分类完成",
		zap.Int64("userId", userID),
		zap.String("category", category))

	// 5. 保存对话历史到 Redis
	s.saveHistory(ctx, userID, question)

	// 6. 路由到对应的 Agent
	go s.routeToAgent(userID, question, category)

	return result, nil
}

// buildClassifyPrompt 构建分类提示词
func (s *ClassifierService) buildClassifyPrompt(question string, history []string) string {
	prompt := "请根据以下问题进行分类：\n\n"
	
	if len(history) > 0 {
		prompt += "对话历史：\n"
		for _, h := range history {
			prompt += h + "\n"
		}
		prompt += "\n"
	}
	
	prompt += "用户问题：" + question + "\n\n"
	prompt += "可选分类：\n"
	for category, info := range s.categories {
		prompt += fmt.Sprintf("- %s: %s\n", category, info.Description)
	}
	
	prompt += "\n请直接返回分类名称，只返回一个词。"
	return prompt
}

// parseCategory 解析分类结果
func (s *ClassifierService) parseCategory(response string) string {
	// 简单匹配（实际应该更智能）
	for category := range s.categories {
		if contains(response, category) {
			return category
		}
	}
	return "general-chat" // 默认分类
}

// routeToAgent 路由到 Agent
func (s *ClassifierService) routeToAgent(userID int64, question string, category string) {
	agentURL, ok := s.agentURLs[category]
	if !ok {
		s.logger.Error("未找到 Agent URL", zap.String("category", category))
		return
	}

	// 调用 Agent API
	reqBody := map[string]interface{}{
		"userId":   userID,
		"question": question,
		"category": category,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		s.logger.Error("序列化请求失败", zap.Error(err))
		return
	}

	resp, err := s.httpClient.Post(agentURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("调用 Agent 失败",
			zap.String("category", category),
			zap.Error(err))
		return
	}
	defer resp.Body.Close()

	s.logger.Info("Agent 调用成功",
		zap.Int64("userId", userID),
		zap.String("category", category))
}

// saveHistory 保存对话历史
func (s *ClassifierService) saveHistory(ctx context.Context, userID int64, message string) {
	historyKey := fmt.Sprintf("chat_history:%d", userID)
	s.redisClient.RPush(ctx, historyKey, message)
	s.redisClient.Expire(ctx, historyKey, 3600*24) // 24 小时过期
}

// 辅助函数
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || 
		len(str) > len(substr) && (str[:len(substr)] == substr || str[len(str)-len(substr):] == substr))
}

