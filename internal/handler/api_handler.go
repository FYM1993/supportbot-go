package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supportbot/supportbot-go/internal/model"
	"github.com/supportbot/supportbot-go/internal/service"
	"go.uber.org/zap"
)

// APIHandler API 处理器
type APIHandler struct {
	sessionService *service.SessionService
	logger         *zap.Logger
}

// NewAPIHandler 创建 API 处理器
func NewAPIHandler(sessionService *service.SessionService, logger *zap.Logger) *APIHandler {
	return &APIHandler{
		sessionService: sessionService,
		logger:         logger,
	}
}

// ReceiveAIResponse 接收 AI 回复
func (h *APIHandler) ReceiveAIResponse(c *gin.Context) {
	var req model.AIResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	h.logger.Info("收到 AI 回复",
		zap.Int64("userId", req.UserID),
		zap.String("source", req.Source),
		zap.String("content", req.Content))

	// 构建消息
	msg := model.ChatMessage{
		MessageID:  uuid.New().String(),
		Type:       "AI_RESPONSE",
		Content:    req.Content,
		Sender:     0,
		SenderName: "AI助手",
		Timestamp:  time.Now(),
	}

	// 推送给用户
	err := h.sessionService.SendMessageToUser(req.UserID, msg)
	if err != nil {
		h.logger.Error("推送 AI 回复失败",
			zap.Int64("userId", req.UserID),
			zap.Error(err))
		c.JSON(500, gin.H{"success": false, "message": "推送失败"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "推送成功"})
}

// Health 健康检查
func (h *APIHandler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":        "UP",
		"service":       c.GetString("service_name"),
		"online_users":  h.sessionService.GetOnlineCount(),
	})
}

// UserLogin 用户登录（简化版）
func (h *APIHandler) UserLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if req.Username == "" {
		c.JSON(400, gin.H{"error": "用户名不能为空"})
		return
	}

	// 简单处理：返回用户信息（生产环境应该验证密码）
	// 使用用户名的哈希作为 uid（简化处理）
	userID := int64(len(req.Username) * 123) // 简单的 ID 生成

	h.logger.Info("用户登录",
		zap.String("username", req.Username),
		zap.Int64("userId", userID))

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"uid":      userID,
			"username": req.Username,
		},
	})
}

