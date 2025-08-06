import { memo } from 'react'
import type { ArtNet } from '@/types/artnet'
import type { DmxHistoryPoint } from '@/stores/artNetStore'
import DmxHistoryChart from './DmxHistoryChart'
import type { SelectedUniverse } from '@/types'

interface SelectedInfoDisplayProps {
  selectedUniverse: SelectedUniverse | null
  selectedChannel: ArtNet.DmxChannel | null
  dmxHistory: DmxHistoryPoint[]
}

const SelectedInfoDisplay = memo(({ selectedUniverse, selectedChannel, dmxHistory }: SelectedInfoDisplayProps) => {
  const dmxValue = dmxHistory.length > 0 ? dmxHistory[dmxHistory.length - 1].value : null

  return (
    <div className="text-sm">
      <div className="flex items-center justify-between py-1">
        <span className="text-left font-bold">Address</span>
        <span className="text-right">{selectedUniverse ? selectedUniverse.address : 'None'}</span>
      </div>
      <div className="flex items-center justify-between py-1">
        <span className="text-left font-bold">Universe ID</span>
        <span className="text-right">{selectedUniverse ? selectedUniverse.universe : 'None'}</span>
      </div>
      <div className="flex items-center justify-between py-1">
        <span className="text-left font-bold">Selected Channel</span>
        <span className="text-right">{selectedChannel !== null ? selectedChannel : 'None'}</span>
      </div>
      <div className="flex items-center justify-between py-1">
        <span className="text-left font-bold">Dmx Value</span>
        <span className="text-right">{dmxValue != null ? dmxValue : 'None'}</span>
      </div>
      <div className="pt-2">
        <DmxHistoryChart history={dmxHistory} maxLength={100} />
      </div>
    </div>
  )
})

export default SelectedInfoDisplay
