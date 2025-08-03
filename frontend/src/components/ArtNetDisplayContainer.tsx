import React from 'react'
import ArtNetDisplay from './ArtNetDisplay'
import { useWebSocket } from '@/contexts/WebSocketContext'

interface ArtNetDisplayContainerProps {
  displayUniverse?: [string, number]
}

const ArtNetDisplayContainer: React.FC<ArtNetDisplayContainerProps> = ({ displayUniverse }) => {
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
  return <ArtNetDisplay displayUniverse={displayUniverse} dmxData={dmxDataForDisplay} />
}

export default ArtNetDisplayContainer
