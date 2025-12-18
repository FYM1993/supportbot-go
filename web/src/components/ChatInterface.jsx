import React, { useEffect, useRef } from 'react'
import { LogOut, Send, Bot, User, Wifi, WifiOff } from 'lucide-react'
import { useChat } from '../contexts/ChatContext'
import MessageList from './MessageList'
import MessageInput from './MessageInput'
import ConnectionStatus from './ConnectionStatus'

const ChatInterface = ({ userInfo, onLogout }) => {
  const { 
    messages, 
    isConnected, 
    isTyping, 
    sendMessage, 
    connect, 
    disconnect 
  } = useChat()

  const messagesEndRef = useRef(null)

  // 自动滚动到底部
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  // 登录后自动连接WebSocket
  useEffect(() => {
    if (userInfo?.uid) {
      connect(userInfo.uid)
    }
  }, [userInfo?.uid, connect])

  const handleLogout = () => {
    disconnect()
    onLogout()
  }

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="bg-white rounded-t-2xl shadow-lg border border-gray-100 p-6 mb-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="w-12 h-12 bg-primary-100 rounded-full flex items-center justify-center">
              <Bot className="w-6 h-6 text-primary-600" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-800">
                SupportBot Assistant
              </h2>
              <p className="text-gray-600 text-sm">
                Welcome, {userInfo?.username}!
              </p>
            </div>
          </div>
          
          <div className="flex items-center space-x-4">
            <ConnectionStatus isConnected={isConnected} />
            <button
              onClick={handleLogout}
              className="btn-secondary flex items-center space-x-2"
            >
              <LogOut className="w-4 h-4" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </div>

      {/* Chat Container */}
      <div className="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
        {/* Messages */}
        <div className="h-96 overflow-y-auto p-6">
          <MessageList messages={messages} />
          {isTyping && (
            <div className="flex items-center space-x-2 text-gray-500 text-sm">
              <div className="flex space-x-1">
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
              </div>
              <span>AI is typing...</span>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>

        {/* Input */}
        <div className="border-t border-gray-100 p-6">
          <MessageInput 
            onSendMessage={sendMessage} 
            disabled={!isConnected}
            placeholder={isConnected ? "Type your message..." : "Connecting..."}
          />
        </div>
      </div>

      {/* Connection Info */}
      <div className="mt-4 text-center text-sm text-gray-500">
        <p>User ID: {userInfo?.uid}</p>
        <p>Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
      </div>
    </div>
  )
}

export default ChatInterface
