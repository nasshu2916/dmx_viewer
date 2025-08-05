import React, { createContext, useContext, useState } from 'react'
import type { ReactNode } from 'react'
import type { ArtNet } from '@/types/artnet'

interface SelectionContextType {
  selectedUniverse: [string, ArtNet.Universe] | undefined
  setSelectedUniverse: React.Dispatch<React.SetStateAction<[string, ArtNet.Universe] | undefined>>
  selectedChannel: ArtNet.DmxChannel | null
  setSelectedChannel: React.Dispatch<React.SetStateAction<ArtNet.DmxChannel | null>>
}

const SelectionContext = createContext<SelectionContextType | undefined>(undefined)

export const SelectionProvider = ({ children }: { children: ReactNode }) => {
  const [selectedUniverse, setSelectedUniverse] = useState<[string, ArtNet.Universe] | undefined>(undefined)
  const [selectedChannel, setSelectedChannel] = useState<ArtNet.DmxChannel | null>(null)

  return (
    <SelectionContext.Provider value={{ selectedUniverse, setSelectedUniverse, selectedChannel, setSelectedChannel }}>
      {children}
    </SelectionContext.Provider>
  )
}

export const useSelection = () => {
  const context = useContext(SelectionContext)
  if (!context) throw new Error('useSelection must be used within a SelectionProvider')
  return context
}
