import { useState, useCallback } from 'react'
import type { ArtNet } from '@/types/artnet'
import type { ServerMessage } from '@/types/websocket'

export interface ArtNetStore {
  // State
  dmxData: Record<string, Record<number, ArtNet.DmxValue[]>>
  serverMessages: ServerMessage[]
  artNetNodes: ArtNet.ArtNetNode[]

  // Actions
  updateDmxData: (address: string, universe: number, data: ArtNet.DmxValue[]) => void
  addServerMessage: (message: ServerMessage) => void
  setServerMessages: (messages: ServerMessage[]) => void
  setArtNetNodes: (nodes: ArtNet.ArtNetNode[]) => void
  clearData: () => void
}

export const useArtNetStore = (): ArtNetStore => {
  const [dmxData, setDmxData] = useState<Record<string, Record<number, ArtNet.DmxValue[]>>>({})
  const [serverMessages, setServerMessages] = useState<ServerMessage[]>([])
  const [artNetNodes, setArtNetNodes] = useState<ArtNet.ArtNetNode[]>([])

  const updateDmxData = useCallback((address: string, universe: number, data: ArtNet.DmxValue[]) => {
    setDmxData(prevData => ({
      ...prevData,
      [address]: {
        ...(prevData[address] || {}),
        [universe]: data,
      },
    }))
  }, [])

  const addServerMessage = useCallback((message: ServerMessage) => {
    setServerMessages(prevMessages => {
      // Prevent duplicate messages based on timestamp
      if (prevMessages.length > 0 && prevMessages[prevMessages.length - 1].Timestamp === message.Timestamp) {
        return prevMessages
      }
      return [...prevMessages, message]
    })
  }, [])

  const setServerMessagesCallback = useCallback((messages: ServerMessage[]) => {
    setServerMessages(messages)
  }, [])

  const setArtNetNodesCallback = useCallback((nodes: ArtNet.ArtNetNode[]) => {
    setArtNetNodes(nodes)
  }, [])

  const clearData = useCallback(() => {
    setDmxData({} as Record<string, Record<number, ArtNet.DmxValue[]>>)
    setServerMessages([])
    setArtNetNodes([])
  }, [])

  return {
    dmxData,
    serverMessages,
    artNetNodes,
    updateDmxData,
    addServerMessage,
    setServerMessages: setServerMessagesCallback,
    setArtNetNodes: setArtNetNodesCallback,
    clearData,
  }
}
