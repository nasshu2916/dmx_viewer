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
  // 横並び数（columns）をContainerで管理
  const [columns, setColumns] = React.useState(16)
  // eslint-disable-next-line no-undef
  const containerRef = React.useRef<HTMLDivElement>(null)

  const calcColumns = React.useCallback(() => {
    // ArtNetDisplayContainerの実際の幅を取得
    const containerWidth = containerRef.current?.clientWidth ?? window.innerWidth
    const cellWidth = 48 // セルの最小幅(px)
    const maxColumns = 32
    const minColumns = 1
    let cols = Math.max(minColumns, Math.min(maxColumns, Math.floor(containerWidth / cellWidth)))
    // 2のN乗に切り捨て
    cols = Math.pow(2, Math.floor(Math.log2(cols)))
    if (cols < minColumns) cols = minColumns
    if (cols > maxColumns) cols = maxColumns
    setColumns(cols)
  }, [])

  React.useEffect(() => {
    calcColumns()
    window.addEventListener('resize', calcColumns)
    return () => window.removeEventListener('resize', calcColumns)
  }, [calcColumns])
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

  // キー移動
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (selectedChannel === null) return
    // columnsはContainerのstateを使う
    const maxChannel = 511
    const row = Math.floor(selectedChannel / columns)
    const col = selectedChannel % columns
    let newRow = row
    let newCol = col
    const rows = Math.ceil(512 / columns)
    switch (e.key) {
      case 'ArrowUp':
        newRow = Math.max(0, row - 1)
        break
      case 'ArrowDown':
        newRow = Math.min(rows - 1, row + 1)
        break
      case 'ArrowLeft':
        newCol = Math.max(0, col - 1)
        break
      case 'ArrowRight':
        newCol = Math.min(columns - 1, col + 1)
        break
      default:
        return
    }
    let newChannel = newRow * columns + newCol
    if (newChannel > maxChannel) newChannel = maxChannel
    if (newChannel !== selectedChannel) {
      e.preventDefault()
      onSelectChannel(newChannel as ArtNet.DmxChannel)
    }
  }

  return (
    <div
      className="outline-none focus:ring-1 focus:ring-dmx-accent"
      ref={containerRef}
      style={{ outline: 'none' }}
      tabIndex={0}
      onKeyDown={handleKeyDown}
    >
      <ArtNetDisplay
        columns={columns}
        displayUniverse={displayUniverse}
        dmxData={dmxDataForDisplay}
        selectedChannel={selectedChannel}
        onSelectChannel={onSelectChannel}
      />
    </div>
  )
}

export default ArtNetDisplayContainer
