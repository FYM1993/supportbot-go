package model

import "time"

// ChatMessage 聊天消息
type ChatMessage struct {
	MessageID  string    `json:"messageId"`
	Type       string    `json:"type"` // CHAT, HEARTBEAT, AI_RESPONSE
	Content    string    `json:"content"`
	Sender     int64     `json:"sender"`
	SenderName string    `json:"senderName,omitempty"`
	SessionID  string    `json:"sessionId,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Message   string `json:"message"`
}

// AIResponseRequest AI 回复请求
type AIResponseRequest struct {
	UserID  int64  `json:"userId"`
	Content string `json:"content"`
	Source  string `json:"source"` // assistant, rag, chat
}

// ClassifyRequest 问题分类请求
type ClassifyRequest struct {
	Question string `json:"question"`
	UID      int64  `json:"uid"`
}

// ClassifyResponse 问题分类响应
type ClassifyResponse struct {
	Category    string  `json:"category"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}
