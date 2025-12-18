package service

import (
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

// ChatService 聊天服务
type ChatService struct {
	questionClassifierURL string
	httpClient            *http.Client
	logger                *zap.Logger
}

// NewChatService 创建聊天服务
func NewChatService(questionClassifierURL string, logger *zap.Logger) *ChatService {
	return &ChatService{
		questionClassifierURL: questionClassifierURL,
		httpClient:            &http.Client{},
		logger:                logger,
	}
}

// HandleUserMessage 处理用户消息
func (s *ChatService) HandleUserMessage(userID int64, content string) error {
	s.logger.Info("处理用户消息",
		zap.Int64("userId", userID),
		zap.String("content", content))

	// 调用问题分类服务
	err := s.callQuestionClassifier(userID, content)
	if err != nil {
		s.logger.Error("调用问题分类服务失败",
			zap.Int64("userId", userID),
			zap.Error(err))
		return err
	}

	return nil
}

// callQuestionClassifier 调用问题分类服务
func (s *ChatService) callQuestionClassifier(userID int64, question string) error {
	apiURL := fmt.Sprintf("%s/api/classify?question=%s&uid=%d",
		s.questionClassifierURL,
		url.QueryEscape(question),
		userID)

	s.logger.Debug("调用问题分类服务", zap.String("url", apiURL))

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("问题分类服务返回错误: %d", resp.StatusCode)
	}

	s.logger.Info("问题分类服务调用成功", zap.Int64("userId", userID))
	return nil
}

