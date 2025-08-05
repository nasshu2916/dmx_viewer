import { useEffect, useRef, useState, useMemo } from 'react'
import { WebSocketService, type WebSocketConfig } from './WebSocketService'
import { MessageRouter, type MessageHandler } from './MessageRouter'
import { useArtNetStore, type ArtNetStore } from '@/stores'
import type { ArtNet } from '@/types/artnet'
import { getUniverse } from '@/service/artnet'

const DefaultSubscribeTopics = ['artnet/dmx_packet', 'artnet/nodes']

export interface WebSocketManager {
  // Connection state
  isConnected: boolean

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
        const universe: ArtNet.Universe = getUniverse(packet)
        const dmxValues: ArtNet.DmxValue[] = Array.from(packet.Data)
        const address: string = packet.SourceIP ?? 'unknown'
        const receivedAt: Date = new Date(Date.now())
        artNetStore.updateDmxData(address, universe, dmxValues, receivedAt)
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
        DefaultSubscribeTopics.forEach(topic => {
          wsService.send({ type: 'subscribe', topic: topic })
        })
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
    artNetStore,
    sendMessage,
    subscribe,
    unsubscribe,
    connect,
    disconnect,
  }
}
