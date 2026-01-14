package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supportbot/supportbot-go/internal/client"
	"github.com/supportbot/supportbot-go/internal/config"
	"github.com/supportbot/supportbot-go/internal/handler"
	"github.com/supportbot/supportbot-go/internal/middleware"
	"github.com/supportbot/supportbot-go/internal/service"
	"github.com/supportbot/supportbot-go/pkg/logger"
	"github.com/supportbot/supportbot-go/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/question-classifier.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	zapLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("question-classifier 服务启动中...")

	// 初始化 Redis
	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		zapLogger.Fatal("连接 Redis 失败", zap.Error(err))
	}

	// 初始化 LLM 客户端
	llmClient := client.NewDashScopeClient(cfg.DashScope.APIKey, cfg.DashScope.Model, zapLogger)

	// 分类配置（从配置文件读取）
	categories := map[string]service.CategoryInfo{
		"product.inquiry": {
			Name:        "商品咨询",
			Description: "查询商品信息、库存、价格等",
			Keywords:    []string{"商品", "库存", "价格", "详情"},
			AgentURL:    cfg.Services.Assistant + "/api/process",
		},
		"order.status": {
			Name:        "订单查询",
			Description: "查询订单状态、物流信息等",
			Keywords:    []string{"订单", "物流", "配送", "快递"},
			AgentURL:    cfg.Services.Assistant + "/api/process",
		},
		"knowledge.query": {
			Name:        "知识库查询",
			Description: "查询产品使用手册、常见问题等",
			Keywords:    []string{"怎么用", "如何", "教程", "说明"},
			AgentURL:    cfg.Services.KnowledgeRAG + "/api/rag",
		},
		"general-chat": {
			Name:        "通用对话",
			Description: "闲聊、问候等一般性对话",
			Keywords:    []string{"你好", "聊天", "谢谢"},
			AgentURL:    cfg.Services.GeneralChat + "/api/chat",
		},
	}

	systemPrompt := `你是一个智能客服问题分类助手。请根据用户的问题，判断属于哪个分类。
可选分类：
- product.inquiry: 商品咨询
- order.status: 订单查询
- knowledge.query: 知识库查询
- general-chat: 通用对话

请只返回分类名称。`

	// 初始化服务
	classifierService := service.NewClassifierService(llmClient, redisClient, categories, systemPrompt, zapLogger)

	// 初始化处理器
	classifierHandler := handler.NewClassifierHandler(classifierService, zapLogger)

	// 初始化路由
	r := gin.Default()
	r.Use(middleware.CORS())

	// API 路由
	r.GET("/api/classify", classifierHandler.Classify)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": cfg.Server.Name,
		})
	})

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("question-classifier 服务启动成功",
		zap.Int("port", cfg.Server.Port))

	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}
