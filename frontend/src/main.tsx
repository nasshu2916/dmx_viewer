import { createRoot } from 'react-dom/client'
import App from './App'
import './App.css'
import { WebSocketProvider } from './contexts/WebSocketContext'

const container = document.getElementById('root')
if (container) {
  const root = createRoot(container)
  root.render(
    <WebSocketProvider>
      <App />
    </WebSocketProvider>
  )
}
