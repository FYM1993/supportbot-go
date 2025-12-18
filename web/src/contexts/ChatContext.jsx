import React, { createContext, useContext, useEffect, useState, useCallback } from 'react'
import { useWebSocket } from '../hooks/useWebSocket'
import { useHeartbeat } from '../hooks/useHeartbeat'

const ChatContext = createContext()

export const useChat = () => {
  const context = useContext(ChatContext)
  if (!context) {
    throw new Error('useChat must be used within a ChatProvider')
  }
  return context
}

export const ChatProvider = ({ children }) => {
  const [messages, setMessages] = useState([])
  const [isConnected, setIsConnected] = useState(false)
  const [isTyping, setIsTyping] = useState(false)
  
  const { 
    connect, 
    disconnect, 
    sendMessage, 
    connectionStatus 
  } = useWebSocket()
  
  const { startHeartbeat, stopHeartbeat } = useHeartbeat()

  // 添加消息到列表
  const addMessage = useCallback((message) => {
    setMessages(prev => [...prev, {
      ...message,
      id: Date.now() + Math.random(),
      timestamp: new Date()
    }])
  }, [])

  // 发送消息
  const sendChatMessage = useCallback(async (content) => {
    if (!content.trim()) return

    const userMessage = {
      type: 'CHAT',
      content: content.trim(),
      sender: 'user',
      timestamp: new Date()
    }

    addMessage(userMessage)
    sendMessage(userMessage)
    setIsTyping(true)
  }, [addMessage, sendMessage])

  // 连接WebSocket
  const connectToChat = useCallback(async (userId) => {
    try {
      await connect(userId)
      setIsConnected(true)
      startHeartbeat(userId)
    } catch (error) {
      console.error('Failed to connect to chat:', error)
      setIsConnected(false)
    }
  }, [connect, startHeartbeat])

  // 断开连接
  const disconnectFromChat = useCallback(() => {
    disconnect()
    setIsConnected(false)
    stopHeartbeat()
    setMessages([])
  }, [disconnect, stopHeartbeat])

  // 监听连接状态变化
  useEffect(() => {
    setIsConnected(connectionStatus === 'connected')
  }, [connectionStatus])

  // 监听WebSocket消息
  useEffect(() => {
    const handleMessage = (message) => {
      if (message.type === 'AI_RESPONSE') {
        setIsTyping(false)
        addMessage({
          type: 'AI_RESPONSE',
          content: message.content,
          sender: 'ai',
          timestamp: new Date()
        })
      } else if (message.type === 'SYSTEM') {
        addMessage({
          type: 'SYSTEM',
          content: message.content,
          sender: 'system',
          timestamp: new Date()
        })
      }
    }

    // 这里需要从useWebSocket hook中获取消息处理函数
    // 暂时使用模拟实现
  }, [addMessage])

  const value = {
    messages,
    isConnected,
    isTyping,
    sendMessage: sendChatMessage,
    connect: connectToChat,
    disconnect: disconnectFromChat,
    connectionStatus
  }

  return (
    <ChatContext.Provider value={value}>
      {children}
    </ChatContext.Provider>
  )
}
