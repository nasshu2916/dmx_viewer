import type { ArtNet } from '@/types/artnet'
import type { DmxHistoryPoint } from '@/hooks/useDmxHistory'
import DmxHistoryChart from './DmxHistoryChart'

interface SelectedInfoDisplayProps {
  selectedUniverse?: [string, ArtNet.Universe]
  selectedChannel: ArtNet.DmxChannel | null
  dmxValue: number | null
  dmxHistory: DmxHistoryPoint[]
}

const SelectedInfoDisplay = ({ selectedUniverse, selectedChannel, dmxValue, dmxHistory }: SelectedInfoDisplayProps) => {
  return (
    <div className="space-y-2 text-sm">
      <div>
        <span className="font-bold">Address: </span>
        {selectedUniverse ? selectedUniverse[0] : 'None'}
      </div>
      <div>
        <span className="font-bold">Universe ID: </span>
        {selectedUniverse ? selectedUniverse[1] : 'None'}
      </div>
      <div>
        <span className="font-bold">Selected Channel: </span>
        {selectedChannel !== null ? selectedChannel : 'None'}
      </div>
      <div>
        <span className="font-bold">Dmx Value: </span>
        {dmxValue !== null && dmxValue !== undefined ? dmxValue : 'None'}
      </div>
      <div className="pt-2">
        <DmxHistoryChart history={dmxHistory} maxLength={100} />
      </div>
    </div>
  )
}

export default SelectedInfoDisplay
