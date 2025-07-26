import React, { createContext, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { logger } from '../utils/logger'
import type { ArtNet } from '@/types/artnet'
import type { ServerMessage } from '@/types/websocket'

interface WebSocketContextType {
  // WebSocket connection state
  isConnected: boolean
  ws: WebSocket | null

  // Data states
  dmxData: Record<number, ArtNet.DmxValue[]>
  serverMessages: ServerMessage[]
  artNetNodes: ArtNet.ArtNetNode[]

  // Methods
  sendMessage: (message: unknown) => void
  subscribe: (topic: string) => void
  unsubscribe: (topic: string) => void
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined)

interface WebSocketProviderProps {
  children: ReactNode
  wsUrl?: string
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({
  children,
  wsUrl = import.meta.env.VITE_WEBSOCKET_URL,
}) => {
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [isConnected, setIsConnected] = useState<boolean>(false)
  const [dmxData, setDmxData] = useState<Record<number, ArtNet.DmxValue[]>>({})
  const [serverMessages, setServerMessages] = useState<ServerMessage[]>([])
  const [artNetNodes, setArtNetNodes] = useState<ArtNet.ArtNetNode[]>([])

  // WebSocket connection management
  useEffect(() => {
    if (!wsUrl) {
      logger.error('WebSocket URL is not provided')
      return
    }

    let websocket: WebSocket | null = null
    let reconnectAttempts = 0
    const maxReconnectAttempts = 30
    const reconnectInterval = 1000

    const connectWebSocket = () => {
      websocket = new WebSocket(wsUrl)

      websocket.onopen = () => {
        logger.info('WebSocket connected')
        setIsConnected(true)
        setWs(websocket)
        reconnectAttempts = 0 // Reset on successful connection

        // Subscribe to default topics
        websocket?.send(JSON.stringify({ type: 'subscribe', topic: 'artnet/packet' }))
        websocket?.send(JSON.stringify({ type: 'subscribe', topic: 'artnet/dmx_packet' }))
        websocket?.send(JSON.stringify({ type: 'subscribe', topic: 'server/message' }))
        websocket?.send(JSON.stringify({ type: 'subscribe', topic: 'server/message_history' }))
        websocket?.send(JSON.stringify({ type: 'subscribe', topic: 'artnet/nodes' }))
      }

      websocket.onerror = error => {
        logger.error('WebSocket error:', error)
        setIsConnected(false)
        websocket?.close()
      }

      websocket.onclose = () => {
        logger.info('WebSocket disconnected')
        setIsConnected(false)
        setWs(null)

        if (reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++
          logger.info(`Attempting to reconnect WebSocket (attempt ${reconnectAttempts}/${maxReconnectAttempts})...`)
          setTimeout(connectWebSocket, reconnectInterval)
        } else {
          logger.error('Max reconnect attempts reached. WebSocket will not reconnect automatically.')
        }
      }

      websocket.onmessage = event => {
        try {
          const data = JSON.parse(event.data)
          if (!data || !data.Type) {
            logger.warn('Received malformed WebSocket message:', data)
            return
          }

          switch (data.Type) {
            case 'artnet_dmx_packet': {
              const artNetPacket: ArtNet.ArtDMXPacket = data.Data
              artNetPacket.Data = data.Data.DataValues || []
              const universe = (artNetPacket.Net << 8) | artNetPacket.SubUni
              setDmxData(prevData => ({
                ...prevData,
                [universe]: Array.from(artNetPacket.Data) as ArtNet.DmxValue[],
              }))
              break
            }

            case 'server_message': {
              setServerMessages(prevMessages => {
                const newMessage: ServerMessage = data
                if (
                  prevMessages.length > 0 &&
                  prevMessages[prevMessages.length - 1].Timestamp === newMessage.Timestamp
                ) {
                  return prevMessages
                }
                return [...prevMessages, newMessage]
              })
              break
            }
            case 'server_message_history': {
              setServerMessages(data.Data)
              break
            }
            case 'artnet_nodes': {
              const nodes: ArtNet.ArtNetNode[] = data.Data
              setArtNetNodes(nodes)
              break
            }

            default:
              logger.info('Received message type:', data.Type, data)
          }
        } catch (e) {
          logger.error('Failed to parse WebSocket message:', e)
        }
      }
    }

    connectWebSocket()

    return () => {
      websocket?.close()
    }
  }, [wsUrl])

  const sendMessage = (message: unknown) => {
    if (ws && isConnected) {
      ws.send(JSON.stringify(message))
    } else {
      logger.warn('WebSocket is not connected. Message not sent:', message)
    }
  }

  const subscribe = (topic: string) => {
    sendMessage({ type: 'subscribe', topic })
  }

  const unsubscribe = (topic: string) => {
    sendMessage({ type: 'unsubscribe', topic })
  }

  const contextValue: WebSocketContextType = {
    isConnected,
    ws,
    dmxData,
    serverMessages,
    artNetNodes,
    sendMessage,
    subscribe,
    unsubscribe,
  }

  return <WebSocketContext.Provider value={contextValue}>{children}</WebSocketContext.Provider>
}

export const useWebSocket = (): WebSocketContextType => {
  const context = useContext(WebSocketContext)
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }
  return context
}
