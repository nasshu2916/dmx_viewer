import { useEffect, useState } from 'react'
import { logger } from '../utils/logger'
import type { ArtNet } from '@/types/artnet'

interface ServerMessage {
  Type: string
  Message: string
  Timestamp: number
}

const useArtNetWebSocket = () => {
  const [dmxData, setDmxData] = useState<Record<number, ArtNet.DmxValue[]>>({})
  const [isConnected, setIsConnected] = useState<boolean>(false)
  const [serverMessages, setServerMessages] = useState<ServerMessage[]>([])
  const [artNetNodes, setArtNetNodes] = useState<ArtNet.ArtNetNode[]>([])

  // Fetch initial DMX data on hook mount
  useEffect(() => {
    const fetchDMXData = async () => {
      try {
        const res = await fetch('/api/dmx_data')
        if (!res.ok) {
          console.error('Failed to fetch initial DMX data')
          return
        }
        const data = await res.json()
        const parsedData: Record<number, ArtNet.DmxValue[]> = {}
        for (const key in data) {
          parsedData[Number(key)] = Array.from(data[key]) as ArtNet.DmxValue[]
        }
        setDmxData(parsedData)
      } catch (e) {
        logger.error('Error fetching initial DMX data:', e)
      }
    }
    fetchDMXData()
  }, [])

  useEffect(() => {
    const wsUrl = import.meta.env.VITE_WEBSOCKET_URL
    let ws: WebSocket | null = null
    let reconnectAttempts = 0
    const maxReconnectAttempts = 30
    const reconnectInterval = 1000

    const connectWebSocket = () => {
      ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        logger.info('WebSocket connected')
        setIsConnected(true)

        ws?.send(JSON.stringify({ type: 'subscribe', topic: 'artnet/packet' }))
        ws?.send(JSON.stringify({ type: 'subscribe', topic: 'artnet/dmx_packet' }))
        ws?.send(JSON.stringify({ type: 'subscribe', topic: 'server/message' }))
        ws?.send(JSON.stringify({ type: 'subscribe', topic: 'server/message_history' }))
        ws?.send(JSON.stringify({ type: 'subscribe', topic: 'artnet/nodes' }))
        // ws?.send(JSON.stringify({ type: 'get_nodes' }))
      }

      ws.onerror = error => {
        logger.error('WebSocket error:', error)
        setIsConnected(false)
        ws?.close() // Ensure the socket is closed on error
      }

      ws.onclose = () => {
        logger.info('WebSocket disconnected')
        setIsConnected(false)
        if (reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++
          logger.info(`Attempting to reconnect WebSocket (attempt ${reconnectAttempts}/${maxReconnectAttempts})...`)
          setTimeout(connectWebSocket, reconnectInterval)
        } else {
          logger.error('Max reconnect attempts reached. WebSocket will not reconnect automatically.')
        }
      }

      ws.onmessage = event => {
        try {
          const data = JSON.parse(event.data)
          if (!data || !data.Type) {
            logger.warn('Received malformed WebSocket message:', data)
            return
          }
          if (data.Type === 'artnet_dmx_packet') {
            const artNetPacket: ArtNet.ArtDMXPacket = data.Data
            artNetPacket.Data = data.Data.DataValues || []
            const universe = (artNetPacket.Net << 8) | artNetPacket.SubUni
            setDmxData(prevData => ({
              ...prevData,
              [universe]: Array.from(artNetPacket.Data) as ArtNet.DmxValue[],
            }))
          } else if (data.Type === 'server_message') {
            setServerMessages(prevMessages => {
              const newMessage: ServerMessage = data
              if (prevMessages.length > 0 && prevMessages[prevMessages.length - 1].Timestamp === newMessage.Timestamp) {
                return prevMessages // Message already exists, do not add
              }
              return [...prevMessages, newMessage]
            })
          } else if (data.Type === 'server_message_history') {
            setServerMessages(data.Data)
          } else if (data.Type === 'artnet_nodes') {
            const nodes: ArtNet.ArtNetNode[] = JSON.parse(data.Data)
            setArtNetNodes(nodes)
          } else {
            console.log('Received other message type:', data.Type)
          }
        } catch (e) {
          logger.error('Failed to parse WebSocket message:', e)
        }
      }
    }

    connectWebSocket()

    return () => {
      ws?.close()
    }
  }, [])

  return { dmxData, isConnected, serverMessages, artNetNodes } // Return artNetNodes
}

export default useArtNetWebSocket
