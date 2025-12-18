package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supportbot/supportbot-go/internal/config"
	"github.com/supportbot/supportbot-go/internal/handler"
	"github.com/supportbot/supportbot-go/internal/middleware"
	"github.com/supportbot/supportbot-go/internal/service"
	"github.com/supportbot/supportbot-go/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/im-demo.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	zapLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("im-demo 服务启动中...")

	// 初始化服务
	sessionService := service.NewSessionService(zapLogger)
	chatService := service.NewChatService(cfg.Services.QuestionClassifier, zapLogger)

	// 初始化处理器
	wsHandler := handler.NewWebSocketHandler(sessionService, chatService, zapLogger)
	apiHandler := handler.NewAPIHandler(sessionService, zapLogger)

	// 初始化路由
	r := gin.Default()
	r.Use(middleware.CORS())

	// WebSocket 端点（原生 WebSocket）
	r.GET("/ws", wsHandler.HandleWebSocket)

	// HTTP API
	r.POST("/api/ai-response/send", apiHandler.ReceiveAIResponse)
	r.POST("/api/user/login", apiHandler.UserLogin)
	r.GET("/api/health", func(c *gin.Context) {
		c.Set("service_name", cfg.Server.Name)
		apiHandler.Health(c)
	})

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("im-demo 服务启动成功",
		zap.Int("port", cfg.Server.Port))

	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}
