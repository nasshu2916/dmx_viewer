import { useEffect, useRef, useState, useMemo } from 'react'
import { WebSocketService, type WebSocketConfig } from './WebSocketService'
import { MessageRouter, type MessageHandler } from './MessageRouter'
import { useDataStore, type DataStore } from './DataStore'
import type { ArtNet } from '@/types/artnet'
import type { ServerMessage } from '@/types/websocket'

export interface WebSocketManager {
  // Connection state
  isConnected: boolean

  // Data store properties (flattened)
  dmxData: Record<number, ArtNet.DmxValue[]>
  serverMessages: ServerMessage[]
  artNetNodes: ArtNet.ArtNetNode[]

  // Data store (for internal use)
  dataStore: DataStore

  // Actions
  sendMessage: (message: unknown) => boolean
  subscribe: (topic: string) => void
  unsubscribe: (topic: string) => void
  connect: () => void
  disconnect: () => void
}

export const useWebSocketManager = (config: WebSocketConfig): WebSocketManager => {
  const [isConnected, setIsConnected] = useState(false)
  const dataStore = useDataStore()
  const wsServiceRef = useRef<WebSocketService | null>(null)
  const messageRouterRef = useRef<MessageRouter | null>(null)

  // Memoize config to prevent unnecessary re-initialization
  const stableConfig = useMemo(() => config, [config.url, config.maxReconnectAttempts, config.reconnectInterval])

  // Initialize services
  useEffect(() => {
    const wsService = new WebSocketService(stableConfig)
    const messageRouter = new MessageRouter()

    wsServiceRef.current = wsService
    messageRouterRef.current = messageRouter

    // Set up message handlers
    const messageHandlers: MessageHandler = {
      onArtNetDmxPacket: (packet: ArtNet.ArtDMXPacket) => {
        const universe = (packet.Net << 8) | packet.SubUni
        const dmxValues = Array.from(packet.Data) as ArtNet.DmxValue[]
        dataStore.updateDmxData(universe, dmxValues)
      },
      onServerMessage: message => {
        dataStore.addServerMessage(message)
      },
      onServerMessageHistory: messages => {
        dataStore.setServerMessages(messages)
      },
      onArtNetNodes: nodes => {
        dataStore.setArtNetNodes(nodes)
      },
    }

    messageRouter.setHandlers(messageHandlers)

    // Set up WebSocket event handlers
    wsService.setEventHandlers({
      onOpen: () => {
        setIsConnected(true)
        // Subscribe to default topics on connection
        wsService.send({ type: 'subscribe', topic: 'artnet/dmx_packet' })
        wsService.send({ type: 'subscribe', topic: 'artnet/nodes' })
      },
      onClose: () => {
        setIsConnected(false)
      },
      onError: () => {
        setIsConnected(false)
      },
      onMessage: data => {
        messageRouter.route(data)
      },
    })

    // Auto-connect
    // wsService.connect()

    return () => {
      wsService.disconnect()
    }
  }, [stableConfig])

  const sendMessage = (message: unknown): boolean => {
    return wsServiceRef.current?.send(message) ?? false
  }

  const subscribe = (topic: string): void => {
    sendMessage({ type: 'subscribe', topic })
  }

  const unsubscribe = (topic: string): void => {
    sendMessage({ type: 'unsubscribe', topic })
  }

  const connect = (): void => {
    wsServiceRef.current?.connect()
  }

  const disconnect = (): void => {
    wsServiceRef.current?.disconnect()
    dataStore.clearData()
  }

  return {
    isConnected,
    dmxData: dataStore.dmxData,
    serverMessages: dataStore.serverMessages,
    artNetNodes: dataStore.artNetNodes,
    dataStore,
    sendMessage,
    subscribe,
    unsubscribe,
    connect,
    disconnect,
  }
}
