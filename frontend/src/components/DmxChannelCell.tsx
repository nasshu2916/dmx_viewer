import React, { memo } from 'react'
import type { ArtNet } from '@/types/artnet'

interface DmxChannelCellProps {
  channel: ArtNet.DmxChannel
  value: ArtNet.DmxValue
  selected: boolean
  onClick: (channel: ArtNet.DmxChannel) => void
}

const DmxChannelCell: React.FC<DmxChannelCellProps> = memo(({ channel, value, selected, onClick }) => {
  const barHeight = `${(value / 255) * 100}%`

  return (
    <div
      className={`relative flex min-h-[3rem] cursor-pointer flex-col justify-between overflow-hidden border-2 px-2 py-1 text-center font-mono ${selected ? 'bg-dmx-channel-selected border-dmx-accent' : 'border-transparent'}`}
      onClick={() => onClick(channel)}
    >
      <div className="absolute bottom-0 left-0 z-0 w-full bg-dmx-channel-active" style={{ height: barHeight }} />
      <div className="z-10">
        <div className="text-xxs text-dmx-text-gray">{channel}</div>
        <div className="text-sm font-bold">{value}</div>
      </div>
    </div>
  )
})

export default DmxChannelCell
