import { useState } from 'react'
import './App.css'
import ArtNetDisplay from './components/ArtNetDisplay'
import TimeDisplay from './components/TimeDisplay'
import WebSocketStatusIndicator from './components/WebSocketStatusIndicator'
import NodeListDisplay from './components/NodeListDisplay'
import { useWebSocket } from '@/contexts/WebSocketContext'

function App() {
  const { dmxData, isConnected, serverMessages, artNetNodes } = useWebSocket()
  const [selectedUniverse, setSelectedUniverse] = useState<[string, number] | undefined>(undefined)

  const handleUniverseSelection = (address: string, universe: number) => {
    setSelectedUniverse([address, universe])
  }

  return (
    <div className="App flex min-h-screen flex-col bg-dmx-dark-bg text-dmx-text-light">
      <header className="App-header flex items-center justify-between bg-dmx-medium-bg p-4 shadow-md">
        <h1 className="text-2xl font-bold text-dmx-text-light">DMX Viewer</h1>
        <div className="flex items-center space-x-4">
          <TimeDisplay />
          <WebSocketStatusIndicator isConnected={isConnected} />
        </div>
      </header>
      <main className="App-main-content flex flex-1 space-x-4 p-4">
        <div className="w-1/4 rounded-lg bg-dmx-medium-bg p-4 shadow-lg">
          <NodeListDisplay artNetNodes={artNetNodes} dmxData={dmxData} onSelectUniverses={handleUniverseSelection} />
        </div>
        <div className="flex-1 rounded-lg bg-dmx-medium-bg p-4 shadow-lg">
          <ArtNetDisplay displayUniverse={selectedUniverse} dmxData={dmxData} />
        </div>
        <div className="w-1/4 rounded-lg bg-dmx-medium-bg p-4 shadow-lg">
          <h3 className="mb-4 text-lg font-bold text-dmx-text-light">Settings</h3>
        </div>
      </main>
      <footer className="flex items-center justify-center bg-dmx-medium-bg p-2 text-sm text-dmx-text-light">
        {serverMessages.length > 0 && (
          <div className="server-message">{serverMessages[serverMessages.length - 1].Message}</div>
        )}
      </footer>
    </div>
  )
}

export default App
