import React from 'react'
import { Wifi, WifiOff } from 'lucide-react'

const ConnectionStatus = ({ isConnected }) => {
  return (
    <div className={`flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium ${
      isConnected 
        ? 'bg-green-100 text-green-700' 
        : 'bg-red-100 text-red-700'
    }`}>
      {isConnected ? (
        <>
          <Wifi className="w-4 h-4" />
          <span>Connected</span>
        </>
      ) : (
        <>
          <WifiOff className="w-4 h-4" />
          <span>Disconnected</span>
        </>
      )}
    </div>
  )
}

export default ConnectionStatus
