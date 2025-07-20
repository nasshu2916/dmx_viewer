import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import App from './App'
import { describe, it, expect, vi } from 'vitest'

// Mock fetch and WebSocket for testing
vi.mock('./utils/logger', () => ({
  logger: {
    info: vi.fn(),
    error: vi.fn(),
    warn: vi.fn(),
  },
}))

vi.mock('node-fetch', () => ({
  default: vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({}), // Mock empty DMX data
    })
  ),
}))

const mockWebSocket = vi.fn(() => ({
  onopen: null,
  onerror: null,
  onmessage: null,
  close: vi.fn(),
}))

vi.stubGlobal('WebSocket', mockWebSocket)

vi.stubGlobal('import', {
  meta: {
    env: {
      VITE_WEBSOCKET_URL: 'ws://localhost:8080/ws',
    },
  },
})

// Appコンポーネントの基本的なレンダリングテスト

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />)
  })
  it('現在時刻が正しく表示される', () => {
    // 2025-07-07T12:34:56+09:00 を固定時刻とする
    const fixedDate = new Date('2025-07-07T12:34:56+09:00')
    vi.setSystemTime(fixedDate)
    render(<App />)
    // フォーマット例: 12:34:56 JST
    expect(screen.getByText(/12:34:56/)).toBeInTheDocument()
    // タイマーを元に戻す
    vi.useRealTimers()
  })
})
