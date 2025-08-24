import './App.css'
import ArtNetDisplayContainer from './components/ArtNetDisplayContainer'
import TimeDisplayContainer from './components/TimeDisplayContainer'
import SelectedInfoDisplay from './components/SelectedInfoDisplay'
import WebSocketStatusIndicator from './components/WebSocketStatusIndicator'
import NodeListDisplayContainer from './components/NodeListDisplayContainer'
import MobileTabs, { type MobileTabKey } from './components/MobileTabs'
import { useState, useCallback } from 'react'
import { useWebSocket } from '@/contexts/WebSocketContext'
import { useSelectionStore } from '@/stores/selectionStore'
import { useArtNetStore } from '@/stores/artNetStore'

function App() {
  const { isConnected } = useWebSocket()
  const { serverMessages } = useArtNetStore()
  const { selectedUniverse, selectedChannel } = useSelectionStore()
  const { dmxHistory } = useArtNetStore()

  const [activeTab, setActiveTab] = useState<MobileTabKey>('viewer')
  const handleTabChange = useCallback((key: MobileTabKey) => {
    setActiveTab(key)
  }, [])

  return (
    <div className="App flex h-screen min-h-screen flex-col bg-dmx-dark-bg text-dmx-text-light">
      <header className="App-header flex items-center justify-between bg-dmx-medium-bg p-1 md:p-2 border-b border-dmx-border">
        <h1 className="text-base md:text-lg leading-none font-semibold text-dmx-text-light">DMX Viewer</h1>
        <div className="flex items-center space-x-1 md:space-x-2 text-[10px] md:text-xs leading-none">
          <TimeDisplayContainer />
          <WebSocketStatusIndicator isConnected={isConnected} />
        </div>
      </header>
      <main className="App-main-content flex h-full min-h-0 flex-1 flex-col gap-3 p-2 md:flex-row md:gap-4 md:p-4">
        {/* Mobile: tabbed single-pane layout */}
        <div className="md:hidden space-y-3">
          <MobileTabs active={activeTab} onChange={handleTabChange} />
          <div className="h-[calc(100vh-11rem)] min-h-0 overflow-auto rounded-md border border-dmx-border bg-dmx-medium-bg p-2">
            {activeTab === 'nodes' && <NodeListDisplayContainer />}
            {activeTab === 'viewer' && <ArtNetDisplayContainer />}
            {activeTab === 'status' && (
              <div>
                <h3 className="mb-3 text-lg font-semibold text-dmx-text-light">Status</h3>
                <SelectedInfoDisplay
                  dmxHistory={dmxHistory}
                  selectedChannel={selectedChannel}
                  selectedUniverse={selectedUniverse}
                />
              </div>
            )}
          </div>
        </div>

        {/* Desktop: three-pane layout */}
        <div className="hidden h-full min-h-0 w-full md:block md:w-1/4 md:max-w-xs">
          <div className="h-full max-h-full min-h-0 overflow-auto rounded-md border border-dmx-border bg-dmx-medium-bg p-3">
            <NodeListDisplayContainer />
          </div>
        </div>
        <div className="hidden h-full min-h-0 w-full md:block md:flex-1">
          <div className="h-full min-h-0 overflow-auto rounded-md border border-dmx-border bg-dmx-medium-bg p-3">
            <ArtNetDisplayContainer />
          </div>
        </div>
        <div className="hidden h-full min-h-0 w-full md:block md:w-1/4 md:max-w-xs">
          <div className="h-full max-h-full min-h-0 overflow-auto rounded-md border border-dmx-border bg-dmx-medium-bg p-3">
            <h3 className="mb-3 text-lg font-semibold text-dmx-text-light">Status</h3>
            <SelectedInfoDisplay
              dmxHistory={dmxHistory}
              selectedChannel={selectedChannel}
              selectedUniverse={selectedUniverse}
            />
          </div>
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
