import { useState, useCallback, useRef, useEffect } from 'react'

export const useWebSocket = () => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected')
  const wsRef = useRef(null)
  const reconnectTimeoutRef = useRef(null)
  const userIdRef = useRef(null)

  const connect = useCallback(async (userId) => {
    try {
      setConnectionStatus('connecting')
      userIdRef.current = userId
      
      // 构建 WebSocket URL
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const host = window.location.hostname
      const wsUrl = `${protocol}//${host}:11005/ws?uid=${userId}`
      
      console.log('Connecting to WebSocket:', wsUrl)
      
      // 创建原生 WebSocket 连接
      const ws = new WebSocket(wsUrl)
      
      // 连接成功回调
      ws.onopen = () => {
        console.log('WebSocket connected successfully')
        setConnectionStatus('connected')
      }
      
      // 接收消息回调
      ws.onmessage = (event) => {
        try {
          const messageData = JSON.parse(event.data)
          console.log('Received message:', messageData)
          
          // 使用全局事件通知消息
          window.dispatchEvent(new CustomEvent('websocket-message', {
            detail: messageData
          }))
        } catch (error) {
          console.error('Failed to parse message:', error)
        }
      }
      
      // 连接错误回调
      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        setConnectionStatus('error')
      }
      
      // 连接断开回调
      ws.onclose = (event) => {
        console.log('WebSocket connection closed:', event.code, event.reason)
        setConnectionStatus('disconnected')
        
        // 非正常关闭时尝试重连
        if (event.code !== 1000 && userIdRef.current) {
          if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current)
          }
          reconnectTimeoutRef.current = setTimeout(() => {
            console.log('Attempting to reconnect...')
            connect(userIdRef.current)
          }, 5000)
        }
      }
      
      wsRef.current = ws
      
    } catch (error) {
      console.error('Failed to connect WebSocket:', error)
      setConnectionStatus('error')
      throw error
    }
  }, [])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    
    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnected')
      wsRef.current = null
    }
    
    userIdRef.current = null
    setConnectionStatus('disconnected')
  }, [])

  const sendMessage = useCallback((message) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      try {
        wsRef.current.send(JSON.stringify(message))
        console.log('Message sent:', message)
      } catch (error) {
        console.error('Failed to send message:', error)
        throw error
      }
    } else {
      throw new Error('WebSocket not connected')
    }
  }, [])

  // 清理函数
  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close(1000, 'Component unmounted')
      }
    }
  }, [])

  return {
    connect,
    disconnect,
    sendMessage,
    connectionStatus
  }
}
