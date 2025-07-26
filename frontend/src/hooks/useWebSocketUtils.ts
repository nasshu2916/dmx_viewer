import { useWebSocket } from '../contexts/WebSocketContext'

export const useWebSocketSender = () => {
  const { sendMessage, subscribe, unsubscribe, isConnected } = useWebSocket()

  const sendArtNetCommand = (command: string, data?: unknown) => {
    sendMessage({
      type: 'artnet_command',
      command,
      data,
    })
  }

  const sendServerCommand = (command: string, data?: unknown) => {
    sendMessage({
      type: 'server_command',
      command,
      data,
    })
  }

  const subscribeToTopic = (topic: string) => {
    subscribe(topic)
  }

  const unsubscribeFromTopic = (topic: string) => {
    unsubscribe(topic)
  }

  return {
    sendMessage,
    sendArtNetCommand,
    sendServerCommand,
    subscribeToTopic,
    unsubscribeFromTopic,
    isConnected,
  }
}

/**
 * WebSocketの受信データを監視するフック
 * 特定のデータタイプにフォーカスしたい場合に使用
 */
export const useWebSocketData = () => {
  const { dmxData, serverMessages, artNetNodes, isConnected } = useWebSocket()

  return {
    dmxData,
    serverMessages,
    artNetNodes,
    isConnected,
  }
}

/**
 * WebSocketの生の接続にアクセスするフック
 */
export const useWebSocketConnection = () => {
  const { ws, isConnected } = useWebSocket()

  return {
    ws,
    isConnected,
  }
}
