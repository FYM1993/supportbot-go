import React, { useState } from 'react'
import LoginForm from './components/LoginForm'
import ChatInterface from './components/ChatInterface'
import { ChatProvider } from './contexts/ChatContext'

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [userInfo, setUserInfo] = useState(null)

  const handleLogin = (userData) => {
    setUserInfo(userData)
    setIsLoggedIn(true)
  }

  const handleLogout = () => {
    setUserInfo(null)
    setIsLoggedIn(false)
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-8">
        <header className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">
            SupportBot
          </h1>
          <p className="text-gray-600 text-lg">
            Professional AI-Powered Customer Service
          </p>
        </header>

        {!isLoggedIn ? (
          <LoginForm onLogin={handleLogin} />
        ) : (
          <ChatProvider>
            <ChatInterface 
              userInfo={userInfo} 
              onLogout={handleLogout} 
            />
          </ChatProvider>
        )}
      </div>
    </div>
  )
}

export default App
