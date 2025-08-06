import { logger } from '@/utils/logger'

export interface WebSocketConfig {
  url: string
  maxReconnectAttempts?: number
  reconnectInterval?: number
}

interface WebSocketEventHandlers {
  onOpen?: () => void
  onClose?: () => void
  onError?: () => void
  onMessage?: (data: unknown) => void
}

export class WebSocketService {
  private ws: WebSocket | null = null
  private config: Required<WebSocketConfig>
  private handlers: WebSocketEventHandlers = {}
  private reconnectAttempts = 0
  private isConnecting = false

  constructor(config: WebSocketConfig) {
    this.config = {
      maxReconnectAttempts: 30,
      reconnectInterval: 1000,
      ...config,
    }
  }

  connect(): void {
    logger.info('WebSocketService: connect called')
    if (this.isConnecting || this.isConnected()) {
      logger.info('WebSocket is already connecting or connected')
      return
    }

    if (!this.config.url) {
      logger.error('WebSocket URL is not configured')
      return
    }

    logger.debug('Attempting to connect to WebSocket:', this.config.url)
    this.isConnecting = true

    try {
      this.ws = new WebSocket(this.config.url)
    } catch (error) {
      logger.error('Failed to create WebSocket:', error)
      this.isConnecting = false
      return
    }

    this.ws.onopen = () => {
      logger.info('WebSocket connected')
      this.isConnecting = false
      this.reconnectAttempts = 0
      this.handlers.onOpen?.()
    }

    this.ws.onerror = error => {
      logger.error('WebSocket connection error:', error)
      this.handlers.onError?.()
      this.ws?.close()
    }

    this.ws.onclose = event => {
      logger.info(`WebSocket disconnected - Code: ${event.code}, Reason: ${event.reason}, WasClean: ${event.wasClean}`)
      this.isConnecting = false
      this.handlers.onClose?.()
      this.attemptReconnect()
    }

    this.ws.onmessage = event => {
      try {
        const data = JSON.parse(event.data)
        this.handlers.onMessage?.(data)
      } catch (error) {
        logger.error('Failed to parse WebSocket message:', error)
      }
    }
  }

  disconnect(): void {
    logger.info('WebSocketService: disconnect called')
    this.reconnectAttempts = this.config.maxReconnectAttempts // Prevent reconnection
    if (this.ws) {
      this.ws.close(1000, 'Normal Closure')
      this.ws = null
    }
  }

  send(message: unknown): boolean {
    if (!this.isConnected()) {
      logger.warn('WebSocket is not connected. Message not sent:', message)
      return false
    }

    try {
      this.ws!.send(JSON.stringify(message))
      return true
    } catch (error) {
      logger.error('Failed to send WebSocket message:', error)
      return false
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  setEventHandlers(handlers: WebSocketEventHandlers): void {
    this.handlers = { ...this.handlers, ...handlers }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
      logger.error('Max reconnect attempts reached. WebSocket will not reconnect automatically.')
      return
    }

    this.reconnectAttempts++
    logger.info(
      `Attempting to reconnect WebSocket (attempt ${this.reconnectAttempts}/${this.config.maxReconnectAttempts})...`
    )

    setTimeout(() => {
      this.connect()
    }, this.config.reconnectInterval)
  }
}
