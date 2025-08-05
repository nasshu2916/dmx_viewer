import { create } from 'zustand'
import type { ArtNet } from '@/types/artnet'
import type { ServerMessage } from '@/types/websocket'

export interface ArtNetStore {
  // State
  dmxData: Record<string, Record<ArtNet.Universe, { data: ArtNet.DmxValue[]; receivedAt: Date }>>
  serverMessages: ServerMessage[]
  artNetNodes: ArtNet.ArtNetNode[]

  // Actions
  updateDmxData: (address: string, universe: ArtNet.Universe, data: ArtNet.DmxValue[], receivedAt: Date) => void
  addServerMessage: (message: ServerMessage) => void
  setServerMessages: (messages: ServerMessage[]) => void
  setArtNetNodes: (nodes: ArtNet.ArtNetNode[]) => void
  clearData: () => void
}

export const useArtNetStore = create<ArtNetStore>(set => ({
  dmxData: {},
  serverMessages: [],
  artNetNodes: [],

  updateDmxData: (address, universe, data, receivedAt) => {
    set(state => ({
      dmxData: {
        ...state.dmxData,
        [address]: {
          ...(state.dmxData[address] || {}),
          [universe]: { data, receivedAt },
        },
      },
    }))
  },

  addServerMessage: message => {
    set(state => {
      const prevMessages = state.serverMessages
      if (prevMessages.length > 0 && prevMessages[prevMessages.length - 1].Timestamp === message.Timestamp) {
        return { serverMessages: prevMessages }
      }
      return { serverMessages: [...prevMessages, message] }
    })
  },

  setServerMessages: messages => {
    set({ serverMessages: messages })
  },

  setArtNetNodes: nodes => {
    set({ artNetNodes: nodes })
  },

  clearData: () => {
    set({
      dmxData: {},
      serverMessages: [],
      artNetNodes: [],
    })
  },
}))
