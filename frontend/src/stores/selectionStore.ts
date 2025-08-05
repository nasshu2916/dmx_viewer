import { create } from 'zustand'
import type { ArtNet } from '@/types/artnet'
import type { SelectedUniverse } from '@/types'

interface SelectionState {
  selectedUniverse: SelectedUniverse | null
  setSelectedUniverse: (u: SelectedUniverse | null) => void
  selectedChannel: ArtNet.DmxChannel | null
  setSelectedChannel: (c: ArtNet.DmxChannel | null) => void
}

export const useSelectionStore = create<SelectionState>(set => ({
  selectedUniverse: null,
  setSelectedUniverse: selectedUniverse => set({ selectedUniverse }),
  selectedChannel: null,
  setSelectedChannel: selectedChannel => set({ selectedChannel }),
}))
