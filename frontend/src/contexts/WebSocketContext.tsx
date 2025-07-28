import React, { createContext, useContext, useEffect, useRef } from 'react'
import type { ReactNode } from 'react'
import { useWebSocketManager, type WebSocketManager } from '@/infrastructure/websocket'
import { logger } from '@/utils/logger'

interface WebSocketManagerProviderProps {
  children: ReactNode
  wsUrl?: string
}

const WebSocketManagerContext = createContext<WebSocketManager | undefined>(undefined)

export const WebSocketProvider: React.FC<WebSocketManagerProviderProps> = ({
  children,
  wsUrl = import.meta.env.VITE_WEBSOCKET_URL,
}) => {
  const webSocketManager = useWebSocketManager({
    url: wsUrl || '',
    maxReconnectAttempts: 30,
    reconnectInterval: 1000,
  })

  const connectAttemptedRef = useRef(false)

  useEffect(() => {
    if (!wsUrl) {
      logger.error('WebSocket URL is not provided')
      return
    }

    // Prevent double connection in React StrictMode
    if (connectAttemptedRef.current) {
      logger.info('WebSocketProvider: Connection already attempted, skipping')
      return
    }

    connectAttemptedRef.current = true

    // Auto-connect when component mounts
    webSocketManager.connect()

    // Cleanup on unmount
    return () => {
      webSocketManager.disconnect()
      connectAttemptedRef.current = false
    }
  }, [wsUrl]) // Remove webSocketManager from deps to prevent recreation

  return <WebSocketManagerContext.Provider value={webSocketManager}>{children}</WebSocketManagerContext.Provider>
}

export const useWebSocket = (): WebSocketManager => {
  const context = useContext(WebSocketManagerContext)
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }
  return context
}
