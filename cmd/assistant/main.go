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

type AssistantRequest struct {
	UserID   int64  `json:"userId"`
	Question string `json:"question"`
	Category string `json:"category"`
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/assistant.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	zapLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("assistant 服务启动中...")

	// 初始化 LLM 客户端
	llmClient := client.NewDashScopeClient(cfg.DashScope.APIKey, cfg.DashScope.Model, zapLogger)

	// 业务服务 URL
	imDemoURL := cfg.Services.IMDemo

	// 初始化路由
	r := gin.Default()
	r.Use(middleware.CORS())

	// 处理 Agent 请求
	r.POST("/api/process", func(c *gin.Context) {
		var req AssistantRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		zapLogger.Info("收到 Assistant 请求",
			zap.Int64("userId", req.UserID),
			zap.String("question", req.Question))

		// 异步处理
		go processRequest(req, llmClient, imDemoURL, zapLogger)

		c.JSON(200, gin.H{"status": "processing"})
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": cfg.Server.Name})
	})

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("assistant 服务启动成功", zap.Int("port", cfg.Server.Port))

	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}

func processRequest(req AssistantRequest, llmClient *client.DashScopeClient, 
	imDemoURL string, logger *zap.Logger) {
	
	// 模拟业务查询结果
	var businessResult string
	switch req.Category {
	case "product.inquiry":
		businessResult = "商品名称：智能手表，价格：¥999.00，库存：充足"
	case "order.status":
		businessResult = "订单编号：20231218001，状态：配送中，预计明天送达"
	default:
		businessResult = "暂无相关信息"
	}

	// 使用 LLM 生成友好回复
	systemPrompt := "你是专业客服助手，根据查询结果给用户友好的回复。"
	prompt := fmt.Sprintf("查询结果：%s\n用户问题：%s\n请回复用户。", businessResult, req.Question)
	
	response, err := llmClient.SimpleChat(systemPrompt, prompt)
	if err != nil {
		logger.Error("LLM 调用失败", zap.Error(err))
		response = businessResult
	}

	// 发送到 im-demo
	sendToIM(req.UserID, response, imDemoURL, logger)
}

func sendToIM(userID int64, content string, imDemoURL string, logger *zap.Logger) {
	aiResp := model.AIResponseRequest{
		UserID:  userID,
		Content: content,
		Source:  "assistant",
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

