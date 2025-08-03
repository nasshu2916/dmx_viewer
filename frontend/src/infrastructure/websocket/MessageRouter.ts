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

  setHandlers(handlers: MessageHandler): void {
    this.handlers = { ...this.handlers, ...handlers }
  }

  route(message: unknown): void {
    if (!this.isValidMessage(message)) {
      logger.warn('Received malformed WebSocket message:', message)
      return
    }

    const { Type, Data } = message as { Type: string; Data: unknown }

    switch (Type) {
      case 'artnet_dmx_packet': {
        const artNetPacket: ArtNet.ArtDMXPacket = Data as ArtNet.ArtDMXPacket
        this.handlers.onArtNetDmxPacket?.(artNetPacket)
        break
      }

      case 'server_message': {
        const serverMessage = Data as ServerMessage
        this.handlers.onServerMessage?.(serverMessage)
        break
      }

      case 'server_message_history': {
        const messages = Data as ServerMessage[]
        this.handlers.onServerMessageHistory?.(messages)
        break
      }

      case 'artnet_nodes': {
        const nodes = Data as ArtNet.ArtNetNode[]
        this.handlers.onArtNetNodes?.(nodes)
        break
      }

      default:
        logger.info('Received message type:', Type, Data)
        this.handlers.onUnknownMessage?.(Type, Data)
    }
  }

  private isValidMessage(message: unknown): message is { Type: string; Data: unknown } {
    return message !== null && typeof message === 'object' && 'Type' in message && typeof message.Type === 'string'
  }
}
