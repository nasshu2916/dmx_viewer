import React from 'react'
import ArtNetDisplay from './ArtNetDisplay'
import { useWebSocket } from '@/contexts/WebSocketContext'

interface ArtNetDisplayContainerProps {
  displayUniverse?: [string, number]
}

const ArtNetDisplayContainer: React.FC<ArtNetDisplayContainerProps> = ({ displayUniverse }) => {
  const { dmxData } = useWebSocket()
  return <ArtNetDisplay displayUniverse={displayUniverse} dmxData={dmxData} />
}

export default ArtNetDisplayContainer
