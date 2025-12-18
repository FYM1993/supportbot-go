package model

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// UserSession 用户会话
type UserSession struct {
	UserID        int64
	Username      string
	Conn          *websocket.Conn
	SessionID     string
	ClientIP      string
	LastHeartbeat time.Time
	MissedBeats   int
	mu            sync.RWMutex // 保护会话字段
}

// UpdateHeartbeat 更新心跳时间
func (s *UserSession) UpdateHeartbeat() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastHeartbeat = time.Now()
	s.MissedBeats = 0
}

// IncrementMissedBeats 增加丢失心跳次数
func (s *UserSession) IncrementMissedBeats() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MissedBeats++
}

// ShouldBeCleaned 判断是否应该清理
func (s *UserSession) ShouldBeCleaned() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.MissedBeats >= 3
}

// WriteMessage 向 WebSocket 写入消息（线程安全）
func (s *UserSession) WriteMessage(message interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Conn.WriteJSON(message)
}
