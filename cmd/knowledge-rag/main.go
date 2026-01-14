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
	"github.com/supportbot/supportbot-go/internal/service"
	"github.com/supportbot/supportbot-go/internal/vectorstore"
	"github.com/supportbot/supportbot-go/pkg/logger"
	"go.uber.org/zap"
)

type RAGRequest struct {
	UserID   int64  `json:"userId"`
	Question string `json:"question"`
}

type AddKnowledgeRequest struct {
	ID       string            `json:"id" binding:"required"`
	Content  string            `json:"content" binding:"required"`
	Metadata map[string]string `json:"metadata"`
}

func main() {
	cfg, err := config.LoadConfig("configs/knowledge-rag.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	zapLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("knowledge-rag 服务启动中...")

	// 初始化 LLM 客户端
	llmClient := client.NewDashScopeClient(cfg.DashScope.APIKey, cfg.DashScope.Model, zapLogger)

	// 初始化 Embedding 客户端
	embeddingClient := client.NewEmbeddingClient(cfg.DashScope.APIKey, zapLogger)

	// 初始化向量存储
	vectorStore := vectorstore.NewMemoryVectorStore(zapLogger)

	// 初始化知识库服务
	knowledgeService := service.NewKnowledgeService(embeddingClient, vectorStore, zapLogger)

	// 加载默认知识库
	if err := knowledgeService.InitDefaultKnowledge(); err != nil {
		log.Fatalf("初始化知识库失败: %v", err)
	}

	imDemoURL := cfg.Services.IMDemo

	r := gin.Default()
	r.Use(middleware.CORS())

	// RAG 检索接口
	r.POST("/api/rag", func(c *gin.Context) {
		var req RAGRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		zapLogger.Info("收到 RAG 请求",
			zap.Int64("userId", req.UserID),
			zap.String("question", req.Question))

		go processRAG(req, llmClient, knowledgeService, imDemoURL, zapLogger)
		c.JSON(200, gin.H{"status": "processing"})
	})

	// 添加知识接口
	r.POST("/api/knowledge", func(c *gin.Context) {
		var req AddKnowledgeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := knowledgeService.AddKnowledge(req.ID, req.Content, req.Metadata); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"id":     req.ID,
		})
	})

	// 查询知识接口
	r.GET("/api/knowledge/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
			return
		}

		results, err := knowledgeService.SearchKnowledge(query, 5, 0.7)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"results": results,
			"count":   len(results),
		})
	})

	// 知识库统计接口
	r.GET("/api/knowledge/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"total_documents": vectorStore.Count(),
			"status":          "ready",
		})
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": cfg.Server.Name})
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("knowledge-rag 服务启动成功",
		zap.Int("port", cfg.Server.Port),
		zap.Int("knowledge_count", vectorStore.Count()))

	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}

func processRAG(req RAGRequest, llmClient *client.DashScopeClient,
	knowledgeService *service.KnowledgeService, imDemoURL string, logger *zap.Logger) {

	logger.Info("开始 RAG 检索", zap.String("question", req.Question))

	// 1. 向量检索知识库
	results, err := knowledgeService.SearchKnowledge(req.Question, 3, 0.7)
	if err != nil {
		logger.Error("知识检索失败", zap.Error(err))
		sendToIM(req.UserID, "抱歉，查询知识库时出错了。", imDemoURL, logger)
		return
	}

	// 2. 构建上下文
	knowledgeContext := knowledgeService.BuildContext(results)
	logger.Info("检索完成",
		zap.Int("results", len(results)),
		zap.Float64("top_score", getTopScore(results)))

	// 3. 使用 LLM 生成回答
	systemPrompt := `你是一个电商领域的智能客服，专门负责知识库问答。
请根据检索到的知识库内容，为用户提供准确、友好的回答。
如果知识库中没有相关信息，请礼貌地告知用户，并建议联系人工客服。

回答要求：
1. 称呼用户为"亲"
2. 语气要亲切、专业
3. 回答要简洁，控制在50字以内
4. 如果有多个相关信息，优先使用相似度最高的
5. 适当使用emoji增加亲和力`

	prompt := fmt.Sprintf("%s\n\n用户问题：%s\n\n请基于上述知识回答用户问题。", knowledgeContext, req.Question)

	response, err := llmClient.SimpleChat(systemPrompt, prompt)
	if err != nil {
		logger.Error("LLM 调用失败", zap.Error(err))
		response = "抱歉，暂时无法处理您的问题，请稍后重试。"
	}

	sendToIM(req.UserID, response, imDemoURL, logger)
}

// getTopScore 获取最高相似度得分
func getTopScore(results []vectorstore.SearchResult) float64 {
	if len(results) == 0 {
		return 0
	}
	return results[0].Score
}

func sendToIM(userID int64, content string, imDemoURL string, logger *zap.Logger) {
	aiResp := model.AIResponseRequest{
		UserID:  userID,
		Content: content,
		Source:  "knowledge-rag",
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
