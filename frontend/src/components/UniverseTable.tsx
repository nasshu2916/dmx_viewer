/* global HTMLTableCellElement, ScrollIntoViewOptions */
import React, { useCallback } from 'react'
import DmxChannelCell from './DmxChannelCell'
import type { ArtNet } from '@/types/artnet'

interface UniverseTableProps {
  universe: number
  data: ArtNet.DmxValue[]
  receivedAt?: Date
  selectedChannel: ArtNet.DmxChannel | null
  onSelectChannel: (channel: ArtNet.DmxChannel) => void
  columns: number
}

const UniverseTable: React.FC<UniverseTableProps> = ({
  universe,
  data,
  receivedAt,
  selectedChannel,
  onSelectChannel,
  columns,
}) => {
  const selectedCellRef = React.useRef<
    null | (HTMLTableCellElement & { scrollIntoView: (options?: ScrollIntoViewOptions) => void })
  >(null)

  React.useEffect(() => {
    // テストやSSR環境で current が null の場合も考慮
    if (selectedCellRef.current && typeof selectedCellRef.current.scrollIntoView === 'function') {
      selectedCellRef.current.scrollIntoView({ block: 'nearest', inline: 'nearest' })
    }
  }, [selectedChannel])

  const length = data.length
  const rows = Math.ceil(length / columns)

  const handleCellClick = useCallback(
    (channel: ArtNet.DmxChannel) => {
      onSelectChannel(channel)
    },
    [onSelectChannel]
  )

  return (
    <div className="mb-2 rounded-lg bg-dmx-light-bg p-2 shadow-lg md:mb-4">
      <h4 className="mb-2 flex text-lg font-bold text-dmx-text-light">
        Universe: {universe}
        {receivedAt && (
          <span className="ml-auto text-sm text-dmx-text-gray">receivedAt: {receivedAt.toLocaleString()}</span>
        )}
      </h4>
      <div className="overflow-x-auto">
        <table className="w-full min-w-[600px] table-fixed border-collapse text-xs text-dmx-text-light md:text-base">
          <tbody>
            {Array.from({ length: rows }).map((_, rowIdx) => {
              const rowStartChannel = rowIdx * columns
              return (
                <tr
                  className={rowIdx % 2 === 0 ? 'bg-dmx-medium-bg' : 'bg-dmx-light-bg'}
                  key={`row-${universe}-${rowStartChannel}`}
                >
                  {Array.from({ length: columns }).map((__, colIdx) => {
                    const channel = (rowStartChannel + colIdx) as ArtNet.DmxChannel
                    if (channel >= length) return <td key={`empty-${universe}-${channel}`} />
                    const value = (data[channel] ?? 0) as ArtNet.DmxValue
                    const isSelected = selectedChannel === channel
                    return (
                      <td
                        className="w-8 min-w-8 p-0 md:w-12 md:min-w-12"
                        key={`channel-${universe}-${channel}`}
                        ref={isSelected ? selectedCellRef : undefined}
                      >
                        <DmxChannelCell
                          channel={channel}
                          selected={isSelected}
                          value={value}
                          onClick={handleCellClick}
                        />
                      </td>
                    )
                  })}
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default UniverseTable
