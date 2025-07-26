export interface ServerMessage {
  Type: string
  Message: string
  Timestamp: number
}

export interface WebSocketMessage {
  type: string
  data?: unknown
  topic?: string
  command?: string
}

export interface WebSocketSubscription {
  type: 'subscribe' | 'unsubscribe'
  topic: string
}

export interface ArtNetCommand {
  type: 'artnet_command'
  command: string
  data?: unknown
}

export interface ServerCommand {
  type: 'server_command'
  command: string
  data?: unknown
}
