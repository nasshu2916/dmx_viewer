import React from 'react'

interface WebSocketStatusIndicatorProps {
  isConnected: boolean
}

const WebSocketStatusIndicator: React.FC<WebSocketStatusIndicatorProps> = React.memo(({ isConnected }) => {
  const statusColor = isConnected ? 'bg-green-500' : 'bg-red-500'
  const statusText = isConnected ? 'Connected' : 'Disconnected'

  return (
    <div className="flex items-center space-x-1">
      <span className={`h-2 w-2 md:h-3 md:w-3 rounded-full ${statusColor}`} />
      <span className="text-xs md:text-sm leading-none text-dmx-text-gray">{`WebSocket: ${statusText}`}</span>
    </div>
  )
})

export default WebSocketStatusIndicator
