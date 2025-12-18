import React, { useState, useRef, useEffect } from 'react'
import { Send, Paperclip, Smile } from 'lucide-react'

const MessageInput = ({ onSendMessage, disabled = false, placeholder = "Type your message..." }) => {
  const [message, setMessage] = useState('')
  const textareaRef = useRef(null)

  // 自动调整文本框高度
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`
    }
  }, [message])

  const handleSubmit = (e) => {
    e.preventDefault()
    if (!message.trim() || disabled) return

    onSendMessage(message.trim())
    setMessage('')
    
    // 重置文本框高度
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
    }
  }

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSubmit(e)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex items-end space-x-3">
      <div className="flex-1 relative">
        <textarea
          ref={textareaRef}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={placeholder}
          disabled={disabled}
          className="input-field resize-none min-h-[44px] max-h-32 overflow-y-auto"
          rows={1}
        />
        
        {/* 工具栏 */}
        <div className="absolute bottom-2 right-2 flex items-center space-x-2">
          <button
            type="button"
            className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
            title="Attach file"
          >
            <Paperclip className="w-4 h-4" />
          </button>
          <button
            type="button"
            className="p-1 text-gray-400 hover:text-gray-600 transition-colors"
            title="Emoji"
          >
            <Smile className="w-4 h-4" />
          </button>
        </div>
      </div>
      
      <button
        type="submit"
        disabled={disabled || !message.trim()}
        className="btn-primary p-3 flex-shrink-0 disabled:opacity-50 disabled:cursor-not-allowed"
        title="Send message"
      >
        <Send className="w-5 h-5" />
      </button>
    </form>
  )
}

export default MessageInput
