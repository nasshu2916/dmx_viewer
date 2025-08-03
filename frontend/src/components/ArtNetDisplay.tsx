import React from 'react'
import DmxChannelCell from './DmxChannelCell'
import type { ArtNet } from '@/types/artnet'

interface ArtNetDisplayProps {
  dmxData: Record<string, Record<number, { data: ArtNet.DmxValue[]; receivedAt: Date }>>
  displayUniverse?: [string, number] | undefined
}

interface UniverseTableProps {
  universe: number
  data: ArtNet.DmxValue[]
  receivedAt?: Date
}

const UniverseTable: React.FC<UniverseTableProps> = ({ universe, data, receivedAt }) => {
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
                        <DmxChannelCell channel={channel} value={value} />
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

const ArtNetDisplay: React.FC<ArtNetDisplayProps> = ({ dmxData, displayUniverse }) => {
  const address = displayUniverse ? displayUniverse[0] : 'Unknown'
  const universe = displayUniverse ? displayUniverse[1] : 0
  const filtered = dmxData[address]?.[universe]

  return (
    <div>
      {filtered === undefined ? (
        <p className="text-dmx-text-light">Waiting for ArtNet data...</p>
      ) : (
        <UniverseTable data={filtered.data} receivedAt={filtered.receivedAt} universe={universe} />
      )}
    </div>
  )
}

export default ArtNetDisplay
