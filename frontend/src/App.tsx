import { useMemo } from 'react'
import './App.css'
import ArtNetDisplayContainer from './components/ArtNetDisplayContainer'
import TimeDisplayContainer from './components/TimeDisplayContainer'
import SelectedInfoDisplay from './components/SelectedInfoDisplay'
import WebSocketStatusIndicator from './components/WebSocketStatusIndicator'
import NodeListDisplayContainer from './components/NodeListDisplayContainer'
import { useWebSocket } from '@/contexts/WebSocketContext'
import { useSelectionStore } from '@/stores/selectionStore'
import { useArtNetStore } from '@/stores/artNetStore'
import { useDmxHistory } from './hooks/useDmxHistory'

function App() {
  const { isConnected } = useWebSocket()
  const { serverMessages, dmxData } = useArtNetStore()
  const { selectedUniverse, selectedChannel } = useSelectionStore()

  const dmxValue = useMemo(() => {
    if (!selectedUniverse || selectedChannel === null) return null
    const addr = selectedUniverse.address
    const univ = selectedUniverse.universe
    return dmxData[addr]?.[univ]?.data[selectedChannel] ?? null
  }, [dmxData, selectedUniverse, selectedChannel])

  // ヒストリー管理
  const selectedKey =
    selectedUniverse && selectedChannel !== null
      ? `${selectedUniverse.address}-${selectedUniverse.universe}-${selectedChannel}`
      : ''

  const dmxHistory = useDmxHistory(dmxValue, selectedKey, 100)
  return (
    <div className="App flex h-screen min-h-screen flex-col bg-dmx-dark-bg text-dmx-text-light">
      <header className="App-header flex items-center justify-between bg-dmx-medium-bg p-4 shadow-md">
        <h1 className="text-2xl font-bold text-dmx-text-light">DMX Viewer</h1>
        <div className="flex items-center space-x-4">
          <TimeDisplayContainer />
          <WebSocketStatusIndicator isConnected={isConnected} />
        </div>
      </header>
      <main className="App-main-content flex h-full min-h-0 flex-1 space-x-4 p-4">
        <div className="h-full max-h-full min-h-0 w-1/4 overflow-auto rounded-lg bg-dmx-medium-bg p-4 shadow-lg">
          <NodeListDisplayContainer />
        </div>
        <div className="h-full min-h-0 flex-1 overflow-auto rounded-lg bg-dmx-medium-bg p-4 shadow-lg">
          <ArtNetDisplayContainer />
        </div>
        <div className="h-full max-h-full min-h-0 w-1/4 overflow-auto rounded-lg bg-dmx-medium-bg p-4 shadow-lg">
          <h3 className="mb-4 text-lg font-bold text-dmx-text-light">Status</h3>
          <SelectedInfoDisplay
            dmxHistory={dmxHistory}
            selectedChannel={selectedChannel}
            selectedUniverse={selectedUniverse}
          />
        </div>
      </main>
      <footer className="flex flex-shrink-0 items-center justify-center bg-dmx-medium-bg p-2 text-sm text-dmx-text-light">
        {serverMessages.length > 0 && (
          <div className="server-message">{serverMessages[serverMessages.length - 1].Message}</div>
        )}
      </footer>
    </div>
  )
}

export default App
