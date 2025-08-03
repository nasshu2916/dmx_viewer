import React from 'react'
import ArtNetDisplay from './ArtNetDisplay'
import { useWebSocket } from '@/contexts/WebSocketContext'

import type { ArtNet } from '@/types/artnet'

interface ArtNetDisplayContainerProps {
  displayUniverse?: [string, number]
  selectedChannel: ArtNet.DmxChannel | null
  onSelectChannel: (channel: ArtNet.DmxChannel) => void
}

const ArtNetDisplayContainer: React.FC<ArtNetDisplayContainerProps> = ({
  displayUniverse,
  selectedChannel,
  onSelectChannel,
}) => {
  const { dmxData } = useWebSocket()
  const dmxDataForDisplay = Object.fromEntries(
    Object.entries(dmxData).map(([address, universes]) => [
      address,
      Object.fromEntries(
        Object.entries(universes).map(([universe, obj]) => [
          Number(universe),
          { data: obj.data, receivedAt: obj.receivedAt },
        ])
      ),
    ])
  )
  return (
    <ArtNetDisplay
      displayUniverse={displayUniverse}
      dmxData={dmxDataForDisplay}
      selectedChannel={selectedChannel}
      onSelectChannel={onSelectChannel}
    />
  )
}

export default ArtNetDisplayContainer
