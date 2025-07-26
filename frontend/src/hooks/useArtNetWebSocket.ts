import { useWebSocket } from '../contexts/WebSocketContext'

const useArtNetWebSocket = () => {
  const { dmxData, isConnected, serverMessages, artNetNodes } = useWebSocket()

  return { dmxData, isConnected, serverMessages, artNetNodes }
}

export default useArtNetWebSocket
