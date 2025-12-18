package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/supportbot/supportbot-go/internal/model"
	"github.com/supportbot/supportbot-go/internal/service"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 生产环境应该检查 Origin 白名单
		return true
	},
}

// WebSocketHandler WebSocket 处理器
type WebSocketHandler struct {
	sessionService *service.SessionService
	chatService    *service.ChatService
	logger         *zap.Logger
}

// NewWebSocketHandler 创建 WebSocket 处理器
func NewWebSocketHandler(sessionService *service.SessionService, chatService *service.ChatService, logger *zap.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		sessionService: sessionService,
		chatService:    chatService,
		logger:         logger,
	}
}

// HandleWebSocket WebSocket 连接入口
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 获取用户 ID
	userIDStr := c.Query("uid")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("WebSocket 升级失败", zap.Error(err))
		return
	}
	defer conn.Close()

	// 注册会话
	sessionID := uuid.New().String()
	clientIP := c.ClientIP()
	username := "用户" + userIDStr
	h.sessionService.RegisterUser(userID, username, conn, sessionID, clientIP)
	defer h.sessionService.RemoveUserBySessionID(sessionID)

	h.logger.Info("WebSocket 连接建立",
		zap.Int64("userId", userID),
		zap.String("sessionId", sessionID))

	// 消息循环
	for {
		var msg model.ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket 读取错误", zap.Error(err))
			}
			break
		}

		// 处理消息
		h.handleMessage(userID, &msg)
	}

	h.logger.Info("WebSocket 连接断开", zap.Int64("userId", userID))
}

// handleMessage 处理用户消息
func (h *WebSocketHandler) handleMessage(userID int64, msg *model.ChatMessage) {
	switch msg.Type {
	case "CHAT":
		// 异步处理聊天消息
		go h.chatService.HandleUserMessage(userID, msg.Content)

		// 立即返回确认消息
		response := model.ChatResponse{
			Success:   true,
			MessageID: msg.MessageID,
			Message:   "消息已收到，正在处理中...",
		}
		h.sessionService.SendMessageToUser(userID, response)

	case "HEARTBEAT":
		// 更新心跳时间
		h.sessionService.UpdateHeartbeat(userID)
		h.logger.Debug("收到心跳", zap.Int64("userId", userID))

	default:
		h.logger.Warn("未知消息类型",
			zap.Int64("userId", userID),
			zap.String("type", msg.Type))
	}
}
