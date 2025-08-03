import React from 'react'
import DmxChannelCell from './DmxChannelCell'
import type { ArtNet } from '@/types/artnet'

interface ArtNetDisplayProps {
  dmxData: Record<string, Record<number, ArtNet.DmxValue[]>>
  displayUniverse?: [string, number] | undefined
}

interface UniverseTableProps {
  universe: number
  data: ArtNet.DmxValue[]
}

const UniverseTable: React.FC<UniverseTableProps> = ({ universe, data }) => {
  return (
    <div className="mb-4 rounded-lg bg-dmx-light-bg p-4 shadow-lg">
      <h4 className="mb-2 text-lg font-bold text-dmx-text-light">Universe: {universe}</h4>
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
  const filteredDmxData = dmxData[address]?.[universe]

  return (
    <div>
      {filteredDmxData === undefined ? (
        <p className="text-dmx-text-light">Waiting for ArtNet data...</p>
      ) : (
        <UniverseTable data={filteredDmxData} universe={universe} />
      )}
    </div>
  )
}

export default ArtNetDisplay
