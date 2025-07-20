import React from 'react'

interface WebSocketStatusIndicatorProps {
  isConnected: boolean
}

const WebSocketStatusIndicator: React.FC<WebSocketStatusIndicatorProps> = ({ isConnected }) => {
  const statusColor = isConnected ? 'bg-green-500' : 'bg-red-500'
  const statusText = isConnected ? 'Connected' : 'Disconnected'

  return (
    <div className="flex items-center space-x-2">
      <span className={`h-3 w-3 rounded-full ${statusColor}`} />
      <span className="text-sm text-dmx-text-gray">{`WebSocket: ${statusText}`}</span>
    </div>
  )
}

export default WebSocketStatusIndicator
