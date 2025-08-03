import { logger } from '@/utils/logger'
import type { ArtNet } from '@/types/artnet'
import type { ServerMessage } from '@/types/websocket'

export interface MessageHandler {
  onArtNetDmxPacket?: (packet: ArtNet.ArtDMXPacket) => void
  onServerMessage?: (message: ServerMessage) => void
  onServerMessageHistory?: (messages: ServerMessage[]) => void
  onArtNetNodes?: (nodes: ArtNet.ArtNetNode[]) => void
  onUnknownMessage?: (type: string, data: unknown) => void
}

export class MessageRouter {
  private handlers: MessageHandler = {}
  // Typeごとのハンドラマップ
  private handlerMap: Record<string, (data: unknown) => void> = {
    artnet_dmx_packet: data => {
      this.handlers.onArtNetDmxPacket?.(data as ArtNet.ArtDMXPacket)
    },
    server_message: data => {
      this.handlers.onServerMessage?.(data as ServerMessage)
    },
    server_message_history: data => {
      this.handlers.onServerMessageHistory?.(data as ServerMessage[])
    },
    artnet_nodes: data => {
      this.handlers.onArtNetNodes?.(data as ArtNet.ArtNetNode[])
    },
  }

  setHandlers(handlers: MessageHandler): void {
    this.handlers = { ...this.handlers, ...handlers }
  }

  route(message: unknown): void {
    if (!this.isValidMessage(message)) {
      logger.warn('Received malformed WebSocket message:', message)
      return
    }

    const { Type, Data } = message as { Type: string; Data: unknown }

    if (Type in this.handlerMap) {
      this.handlerMap[Type](Data)
    } else {
      logger.info('Received message type:', Type, Data)
      this.handlers.onUnknownMessage?.(Type, Data)
    }
  }

  private isValidMessage(message: unknown): message is { Type: string; Data: unknown } {
    return message !== null && typeof message === 'object' && 'Type' in message && typeof message.Type === 'string'
  }
}
