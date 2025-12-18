package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supportbot/supportbot-go/internal/service"
	"go.uber.org/zap"
)

// ClassifierHandler 分类处理器
type ClassifierHandler struct {
	classifierService *service.ClassifierService
	logger            *zap.Logger
}

// NewClassifierHandler 创建分类处理器
func NewClassifierHandler(classifierService *service.ClassifierService, logger *zap.Logger) *ClassifierHandler {
	return &ClassifierHandler{
		classifierService: classifierService,
		logger:            logger,
	}
}

// Classify 问题分类接口
func (h *ClassifierHandler) Classify(c *gin.Context) {
	question := c.Query("question")
	uidStr := c.Query("uid")

	if question == "" || uidStr == "" {
		c.JSON(400, gin.H{"error": "question 和 uid 参数不能为空"})
		return
	}

	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}

	h.logger.Info("收到分类请求",
		zap.Int64("uid", uid),
		zap.String("question", question))

	// 分类并路由
	result, err := h.classifierService.ClassifyAndRoute(uid, question)
	if err != nil {
		h.logger.Error("分类失败", zap.Error(err))
		c.JSON(500, gin.H{"error": "分类失败"})
		return
	}

	c.JSON(200, result)
}

