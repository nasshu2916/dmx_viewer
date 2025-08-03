import React from 'react'
import UniverseTable from './UniverseTable'
import { useWebSocket } from '@/contexts/WebSocketContext'
import { calcColumns, getNextChannelByKey } from './artnetDisplayUtils'

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

  const calcColumnsCallback = React.useCallback(() => {
    const containerWidth = containerRef.current?.clientWidth ?? window.innerWidth
    const cols = calcColumns(containerWidth)
    setColumns(cols)
  }, [])

  React.useEffect(() => {
    calcColumnsCallback()
    window.addEventListener('resize', calcColumnsCallback)
    return () => window.removeEventListener('resize', calcColumnsCallback)
  }, [calcColumnsCallback])
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

  const address = displayUniverse ? displayUniverse[0] : 'Unknown'
  const universe = displayUniverse ? displayUniverse[1] : 0
  const filtered = dmxDataForDisplay[address]?.[universe]

  const maxChannel = filtered ? filtered.data.length - 1 : 0

  // キー移動
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (selectedChannel === null) return
    const newChannel = getNextChannelByKey(e.key, selectedChannel, columns, maxChannel)
    if (newChannel !== null) {
      e.preventDefault()
      onSelectChannel(newChannel)
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
      {filtered === undefined ? (
        <p className="text-dmx-text-light">Waiting for ArtNet data...</p>
      ) : (
        <UniverseTable
          columns={columns}
          data={filtered.data}
          receivedAt={filtered.receivedAt}
          selectedChannel={selectedChannel}
          universe={universe}
          onSelectChannel={onSelectChannel}
        />
      )}
    </div>
  )
}

export default ArtNetDisplayContainer
