import React from 'react'
import DmxChannelCell from './DmxChannelCell'
import type { ArtNet } from '@/types/artnet'

interface ArtNetDisplayProps {
  dmxData: Record<string, Record<number, { data: ArtNet.DmxValue[]; receivedAt: Date }>>
  displayUniverse?: [string, number] | undefined
  selectedChannel: ArtNet.DmxChannel | null
  onSelectChannel: (channel: ArtNet.DmxChannel) => void
}

interface UniverseTableProps {
  universe: number
  data: ArtNet.DmxValue[]
  receivedAt?: Date
}

interface UniverseTableSelectableProps extends UniverseTableProps {
  selectedChannel: ArtNet.DmxChannel | null
  onSelectChannel: (channel: ArtNet.DmxChannel) => void
}

const UniverseTable: React.FC<UniverseTableSelectableProps> = ({
  universe,
  data,
  receivedAt,
  selectedChannel,
  onSelectChannel,
}) => {
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
            {Array.from({ length: 32 })
              .map((_, rowIdx) => ({
                id: `row-${rowIdx}`,
                rowIdx: rowIdx,
              }))
              .map(row => (
                <tr className={row.rowIdx % 2 === 0 ? 'bg-dmx-medium-bg' : 'bg-dmx-light-bg'} key={row.id}>
                  {Array.from({ length: 16 }).map((__, colIdx) => {
                    const channel = (row.rowIdx * 16 + colIdx) as ArtNet.DmxChannel
                    const value = (data[channel] ?? 0) as ArtNet.DmxValue
                    return (
                      <td key={`channel-${channel}`}>
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
              ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

const ArtNetDisplay: React.FC<ArtNetDisplayProps> = ({
  dmxData,
  displayUniverse,
  selectedChannel,
  onSelectChannel,
}) => {
  const address = displayUniverse ? displayUniverse[0] : 'Unknown'
  const universe = displayUniverse ? displayUniverse[1] : 0
  const filtered = dmxData[address]?.[universe]

  return (
    <div>
      {filtered === undefined ? (
        <p className="text-dmx-text-light">Waiting for ArtNet data...</p>
      ) : (
        <UniverseTable
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

export default ArtNetDisplay
