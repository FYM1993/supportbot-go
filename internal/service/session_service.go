package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/supportbot/supportbot-go/internal/model"
	"go.uber.org/zap"
)

var (
	ErrUserOffline = fmt.Errorf("用户不在线")
)

// SessionService 会话管理服务
type SessionService struct {
	userSessions  map[int64]*model.UserSession // userId -> session
	sessionToUser map[string]int64             // sessionId -> userId
	mu            sync.RWMutex                 // 读写锁保护
	logger        *zap.Logger
}

// NewSessionService 创建会话管理服务
func NewSessionService(logger *zap.Logger) *SessionService {
	s := &SessionService{
		userSessions:  make(map[int64]*model.UserSession),
		sessionToUser: make(map[string]int64),
		logger:        logger,
	}

	// 启动心跳检测
	go s.heartbeatChecker()

	return s
}

// RegisterUser 注册用户会话
func (s *SessionService) RegisterUser(userID int64, username string, conn *websocket.Conn, sessionID string, clientIP string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 清理旧会话
	if existingSession, ok := s.userSessions[userID]; ok {
		s.logger.Info("用户重新连接，关闭旧连接",
			zap.Int64("userId", userID),
			zap.String("oldSessionId", existingSession.SessionID))
		existingSession.Conn.Close()
		delete(s.sessionToUser, existingSession.SessionID)
	}

	// 创建新会话
	session := &model.UserSession{
		UserID:        userID,
		Username:      username,
		Conn:          conn,
		SessionID:     sessionID,
		ClientIP:      clientIP,
		LastHeartbeat: time.Now(),
		MissedBeats:   0,
	}

	s.userSessions[userID] = session
	s.sessionToUser[sessionID] = userID

	s.logger.Info("用户会话注册成功",
		zap.Int64("userId", userID),
		zap.String("username", username),
		zap.String("sessionId", sessionID))
}

// SendMessageToUser 向指定用户发送消息
func (s *SessionService) SendMessageToUser(userID int64, message interface{}) error {
	s.mu.RLock()
	session, ok := s.userSessions[userID]
	s.mu.RUnlock()

	if !ok {
		s.logger.Warn("用户不在线，消息发送失败", zap.Int64("userId", userID))
		return ErrUserOffline
	}

	// WebSocket 写入需要加锁（通过 session 自己的方法）
	err := session.WriteMessage(message)
	if err != nil {
		s.logger.Error("消息发送失败",
			zap.Int64("userId", userID),
			zap.Error(err))
		// 异步清理无效连接
		go s.RemoveUserByID(userID)
		return err
	}

	s.logger.Info("消息发送成功", zap.Int64("userId", userID))
	return nil
}

// UpdateHeartbeat 更新心跳时间
func (s *SessionService) UpdateHeartbeat(userID int64) bool {
	s.mu.RLock()
	session, ok := s.userSessions[userID]
	s.mu.RUnlock()

	if !ok {
		return false
	}

	session.UpdateHeartbeat()
	s.logger.Debug("心跳已更新", zap.Int64("userId", userID))
	return true
}

// RemoveUserBySessionID 根据 sessionId 移除会话
func (s *SessionService) RemoveUserBySessionID(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userID, ok := s.sessionToUser[sessionID]; ok {
		delete(s.userSessions, userID)
		delete(s.sessionToUser, sessionID)
		s.logger.Info("用户会话已移除",
			zap.Int64("userId", userID),
			zap.String("sessionId", sessionID))
	}
}

// RemoveUserByID 根据 userId 移除会话
func (s *SessionService) RemoveUserByID(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, ok := s.userSessions[userID]; ok {
		delete(s.sessionToUser, session.SessionID)
		delete(s.userSessions, userID)
		s.logger.Info("用户会话已移除", zap.Int64("userId", userID))
	}
}

// GetOnlineCount 获取在线用户数
func (s *SessionService) GetOnlineCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.userSessions)
}

// heartbeatChecker 心跳检测器
func (s *SessionService) heartbeatChecker() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		now := time.Now()
		for userID, session := range s.userSessions {
			timeSinceHeartbeat := now.Sub(session.LastHeartbeat)

			if timeSinceHeartbeat > 60*time.Second {
				session.IncrementMissedBeats()

				if session.ShouldBeCleaned() {
					s.logger.Info("清理无效会话",
						zap.Int64("userId", userID),
						zap.Int("missedBeats", session.MissedBeats))

					session.Conn.Close()
					delete(s.userSessions, userID)
					delete(s.sessionToUser, session.SessionID)
				} else {
					s.logger.Warn("用户心跳丢失",
						zap.Int64("userId", userID),
						zap.Int("missedBeats", session.MissedBeats))
				}
			}
		}

		s.mu.Unlock()
	}
}

