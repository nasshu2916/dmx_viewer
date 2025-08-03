import React from 'react'
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
  const length = data.length
  const rows = Math.ceil(length / columns)

  return (
    <div className="mb-4 rounded-lg bg-dmx-light-bg p-4 shadow-lg">
      <h4 className="mb-2 flex text-lg font-bold text-dmx-text-light">
        Universe: {universe}
        {receivedAt && (
          <span className="ml-auto text-sm text-dmx-text-gray">receivedAt: {receivedAt.toLocaleString()}</span>
        )}
      </h4>
      <div className="overflow-x-auto">
        <table className="w-full min-w-full table-fixed border-collapse text-dmx-text-light">
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
                    return (
                      <td key={`channel-${universe}-${channel}`}>
                        <DmxChannelCell
                          channel={channel}
                          selected={selectedChannel === channel}
                          value={value}
                          onClick={() => onSelectChannel(channel)}
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
