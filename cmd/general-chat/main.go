package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supportbot/supportbot-go/internal/client"
	"github.com/supportbot/supportbot-go/internal/config"
	"github.com/supportbot/supportbot-go/internal/middleware"
	"github.com/supportbot/supportbot-go/internal/model"
	"github.com/supportbot/supportbot-go/pkg/logger"
	"go.uber.org/zap"
)

type ChatRequest struct {
	UserID   int64  `json:"userId"`
	Question string `json:"question"`
}

func main() {
	cfg, err := config.LoadConfig("configs/general-chat.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	zapLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("general-chat 服务启动中...")

	llmClient := client.NewDashScopeClient(cfg.DashScope.APIKey, cfg.DashScope.Model, zapLogger)
	imDemoURL := cfg.Services.IMDemo

	r := gin.Default()
	r.Use(middleware.CORS())

	r.POST("/api/chat", func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		zapLogger.Info("收到通用对话请求",
			zap.Int64("userId", req.UserID),
			zap.String("question", req.Question))

		go processChat(req, llmClient, imDemoURL, zapLogger)
		c.JSON(200, gin.H{"status": "processing"})
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": cfg.Server.Name})
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("general-chat 服务启动成功", zap.Int("port", cfg.Server.Port))

	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}

func processChat(req ChatRequest, llmClient *client.DashScopeClient, 
	imDemoURL string, logger *zap.Logger) {
	
	systemPrompt := "你是一个友好的客服助手，负责与用户进行日常对话。"
	response, err := llmClient.SimpleChat(systemPrompt, req.Question)
	if err != nil {
		logger.Error("LLM 调用失败", zap.Error(err))
		response = "抱歉，我现在有点忙，请稍后再试。"
	}

	sendToIM(req.UserID, response, imDemoURL, logger)
}

func sendToIM(userID int64, content string, imDemoURL string, logger *zap.Logger) {
	aiResp := model.AIResponseRequest{
		UserID:  userID,
		Content: content,
		Source:  "general-chat",
	}

	jsonData, _ := json.Marshal(aiResp)
	resp, err := http.Post(imDemoURL+"/api/ai-response/send", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("发送到 im-demo 失败", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	logger.Info("回复已发送", zap.Int64("userId", userID))
}

