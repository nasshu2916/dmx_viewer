import { useEffect, useRef, useState, useMemo } from 'react'
import { WebSocketService, type WebSocketConfig } from './WebSocketService'
import { MessageRouter, type MessageHandler } from './MessageRouter'
import { useArtNetStore, type ArtNetStore } from '@/stores'
import type { ArtNet } from '@/types/artnet'
import type { ServerMessage } from '@/types/websocket'

export interface WebSocketManager {
  // Connection state
  isConnected: boolean

  // Data store properties (flattened)
  dmxData: Record<number, ArtNet.DmxValue[]>
  serverMessages: ServerMessage[]
  artNetNodes: ArtNet.ArtNetNode[]

  // Art-Net data store (for internal use)
  artNetStore: ArtNetStore

  // Actions
  sendMessage: (message: unknown) => boolean
  subscribe: (topic: string) => void
  unsubscribe: (topic: string) => void
  connect: () => void
  disconnect: () => void
}

export const useWebSocketManager = (config: WebSocketConfig): WebSocketManager => {
  const [isConnected, setIsConnected] = useState(false)
  const artNetStore = useArtNetStore()
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
        artNetStore.updateDmxData(universe, dmxValues)
      },
      onServerMessage: message => {
        artNetStore.addServerMessage(message)
      },
      onServerMessageHistory: messages => {
        artNetStore.setServerMessages(messages)
      },
      onArtNetNodes: nodes => {
        artNetStore.setArtNetNodes(nodes)
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
    artNetStore.clearData()
  }

  return {
    isConnected,
    dmxData: artNetStore.dmxData,
    serverMessages: artNetStore.serverMessages,
    artNetNodes: artNetStore.artNetNodes,
    artNetStore,
    sendMessage,
    subscribe,
    unsubscribe,
    connect,
    disconnect,
  }
}
