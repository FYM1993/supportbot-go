import React from 'react'
import { Bot, User, Info } from 'lucide-react'

const MessageList = ({ messages }) => {
  if (messages.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <Bot className="w-8 h-8 text-gray-400" />
        </div>
        <h3 className="text-lg font-medium text-gray-600 mb-2">
          Welcome to SupportBot!
        </h3>
        <p className="text-gray-500 text-sm">
          Start a conversation by typing a message below.
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {messages.map((message) => (
        <MessageItem key={message.id} message={message} />
      ))}
    </div>
  )
}

const MessageItem = ({ message }) => {
  const getMessageIcon = () => {
    switch (message.sender) {
      case 'user':
        return <User className="w-5 h-5 text-white" />
      case 'ai':
        return <Bot className="w-5 h-5 text-primary-600" />
      case 'system':
        return <Info className="w-5 h-5 text-gray-500" />
      default:
        return null
    }
  }

  const getMessageClasses = () => {
    switch (message.sender) {
      case 'user':
        return 'chat-bubble chat-bubble-user'
      case 'ai':
        return 'chat-bubble chat-bubble-ai'
      case 'system':
        return 'chat-bubble chat-bubble-system'
      default:
        return 'chat-bubble chat-bubble-ai'
    }
  }

  const getMessageAlignment = () => {
    switch (message.sender) {
      case 'user':
        return 'justify-end'
      case 'system':
        return 'justify-center'
      default:
        return 'justify-start'
    }
  }

  const formatTime = (timestamp) => {
    if (!timestamp) return ''
    const date = new Date(timestamp)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <div className={`flex ${getMessageAlignment()}`}>
      <div className={`flex items-start space-x-3 ${message.sender === 'user' ? 'flex-row-reverse space-x-reverse' : ''}`}>
        {message.sender !== 'system' && (
          <div className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 ${
            message.sender === 'user' 
              ? 'bg-primary-500' 
              : 'bg-gray-100'
          }`}>
            {getMessageIcon()}
          </div>
        )}
        
        <div className={`${getMessageClasses()} ${message.sender === 'system' ? 'max-w-md' : ''}`}>
          <p className="text-sm leading-relaxed">{message.content}</p>
          {message.sender !== 'system' && (
            <p className="text-xs text-gray-500 mt-1">
              {formatTime(message.timestamp)}
            </p>
          )}
        </div>
      </div>
    </div>
  )
}

export default MessageList
