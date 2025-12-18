import { useRef, useCallback } from 'react'
import { useWebSocket } from './useWebSocket'

export const useHeartbeat = () => {
  const heartbeatIntervalRef = useRef(null)
  const { sendMessage } = useWebSocket()

  const startHeartbeat = useCallback((userId) => {
    // 停止之前的心跳
    stopHeartbeat()
    
    // 启动新的心跳，每30秒发送一次
    heartbeatIntervalRef.current = setInterval(() => {
      try {
        const heartbeatMessage = {
          userId: userId,
          clientIp: '127.0.0.1', // 这里应该获取真实的客户端IP
          clientPort: window.location.port || 3000
        }
        
        // 发送心跳消息
        if (sendMessage) {
          sendMessage({
            type: 'HEARTBEAT',
            ...heartbeatMessage
          })
        }
        
        console.log('Heartbeat sent:', heartbeatMessage)
      } catch (error) {
        console.error('Failed to send heartbeat:', error)
      }
    }, 30000) // 30秒
    
    console.log('Heartbeat started for user:', userId)
  }, [sendMessage])

  const stopHeartbeat = useCallback(() => {
    if (heartbeatIntervalRef.current) {
      clearInterval(heartbeatIntervalRef.current)
      heartbeatIntervalRef.current = null
      console.log('Heartbeat stopped')
    }
  }, [])

  return {
    startHeartbeat,
    stopHeartbeat
  }
}
