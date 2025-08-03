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
  // フォーカス管理用ref
  const focusRef = React.useRef(null)

  // キーハンドラ
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (selectedChannel === null) return
    const row = Math.floor(selectedChannel / 16)
    const col = selectedChannel % 16
    let newRow = row
    let newCol = col
    switch (e.key) {
      case 'ArrowUp':
        newRow = Math.max(0, row - 1)
        break
      case 'ArrowDown':
        newRow = Math.min(31, row + 1)
        break
      case 'ArrowLeft':
        newCol = Math.max(0, col - 1)
        break
      case 'ArrowRight':
        newCol = Math.min(15, col + 1)
        break
      default:
        return
    }
    const newChannel = (newRow * 16 + newCol) as ArtNet.DmxChannel
    if (newChannel !== selectedChannel) {
      e.preventDefault()
      onSelectChannel(newChannel)
    }
  }

  return (
    <div
      className="outline-none focus:ring-1 focus:ring-dmx-accent"
      ref={focusRef}
      style={{ outline: 'none' }}
      tabIndex={0}
      onKeyDown={handleKeyDown}
    >
      <ArtNetDisplay
        displayUniverse={displayUniverse}
        dmxData={dmxDataForDisplay}
        selectedChannel={selectedChannel}
        onSelectChannel={onSelectChannel}
      />
    </div>
  )
}

export default ArtNetDisplayContainer
