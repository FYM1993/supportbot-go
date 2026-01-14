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
	"github.com/supportbot/supportbot-go/internal/tools"
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

	// 初始化工具注册中心
	toolRegistry := tools.NewRegistry(zapLogger)
	
	// 注册内置工具
	if err := tools.RegisterBuiltinTools(toolRegistry, zapLogger); err != nil {
		log.Fatalf("注册工具失败: %v", err)
	}

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
			zap.String("question", req.Question),
			zap.String("category", req.Category))

		// 异步处理
		go processRequest(req, llmClient, toolRegistry, imDemoURL, zapLogger)

		c.JSON(200, gin.H{"status": "processing"})
	})

	// 查询可用工具
	r.GET("/api/tools", func(c *gin.Context) {
		toolList := toolRegistry.List()
		c.JSON(200, gin.H{
			"tools": toolList,
			"count": len(toolList),
		})
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": cfg.Server.Name})
	})

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("assistant 服务启动成功", 
		zap.Int("port", cfg.Server.Port),
		zap.Int("tools", toolRegistry.Count()))

	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}

func processRequest(req AssistantRequest, llmClient *client.DashScopeClient, 
	toolRegistry *tools.Registry, imDemoURL string, logger *zap.Logger) {
	
	logger.Info("开始处理请求（支持工具调用）",
		zap.String("question", req.Question),
		zap.String("category", req.Category))

	// 构建系统提示词
	systemPrompt := `你是一个专业的电商客服助手，能够查询商品信息、订单状态、物流信息等。

你的能力：
1. 查询商品详情（名称、价格、库存、规格）
2. 查询订单状态（订单进度、支付状态、发货信息）
3. 查询物流信息（当前位置、预计送达时间）
4. 查询商品库存和配送信息

当用户询问商品、订单、物流相关问题时，请使用工具查询实际数据，然后用亲切、专业的语气回答用户。

回答要求：
1. 称呼用户为"亲"
2. 语气要友好、专业
3. 回答要简洁，不要太长
4. 适当使用emoji增加亲和力`

	// 构建对话消息
	messages := []client.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: req.Question},
	}

	// 获取工具定义
	toolDefs := toolRegistry.GetFunctionDefs()
	
	// 调用 LLM（支持工具调用）
	maxIterations := 5 // 最多迭代 5 次（防止无限循环）
	var finalResponse string
	
	for i := 0; i < maxIterations; i++ {
		logger.Info("LLM 调用", zap.Int("iteration", i+1))
		
		resp, err := llmClient.ChatWithTools(messages, toolDefs)
		if err != nil {
			logger.Error("LLM 调用失败", zap.Error(err))
			finalResponse = "抱歉，处理您的请求时出错了。"
			break
		}

		// 检查是否有工具调用
		if len(resp.Output.Choices) > 0 && len(resp.Output.Choices[0].Message.ToolCalls) > 0 {
			// LLM 要求调用工具
			assistantMsg := resp.Output.Choices[0].Message
			logger.Info("LLM 请求调用工具", 
				zap.Int("toolCount", len(assistantMsg.ToolCalls)))
			
			// 将 assistant 消息加入对话历史
			messages = append(messages, assistantMsg)
			
			// 执行所有工具调用
			for _, toolCall := range assistantMsg.ToolCalls {
				logger.Info("执行工具", 
					zap.String("tool", toolCall.Function.Name),
					zap.String("args", toolCall.Function.Arguments))
				
				// 转换为工具系统的格式
				tc := tools.ToolCall{
					ID:   toolCall.ID,
					Type: toolCall.Type,
					Function: tools.ToolCallFunction{
						Name:      toolCall.Function.Name,
						Arguments: toolCall.Function.Arguments,
					},
				}
				
				// 执行工具
				result, err := toolRegistry.Execute(tc)
				if err != nil {
					logger.Error("工具执行失败", 
						zap.String("tool", toolCall.Function.Name),
						zap.Error(err))
					result = map[string]interface{}{
						"error": err.Error(),
					}
				}
				
				// 将工具结果序列化为 JSON
				resultJSON, _ := json.Marshal(result)
				
				// 将工具结果加入对话历史
				messages = append(messages, client.Message{
					Role:       "tool",
					Content:    string(resultJSON),
					ToolCallID: toolCall.ID,
				})
			}
			
			// 继续下一轮，让 LLM 基于工具结果生成回答
			continue
		}
		
		// 没有工具调用，获取最终回答
		if resp.Output.Text != "" {
			finalResponse = resp.Output.Text
		} else if len(resp.Output.Choices) > 0 {
			finalResponse = resp.Output.Choices[0].Message.Content
		} else {
			finalResponse = "抱歉，没有获取到回答。"
		}
		
		logger.Info("获得最终回答", zap.String("response", finalResponse))
		break
	}

	// 发送到 im-demo
	sendToIM(req.UserID, finalResponse, imDemoURL, logger)
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

