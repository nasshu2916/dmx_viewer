import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { WebSocketService } from './WebSocketService'

class MockWebSocket {
  static instances: MockWebSocket[] = []
  static OPEN = 1
  static CLOSED = 3
  public readyState = MockWebSocket.OPEN
  public onopen: (() => void) | null = null
  public onclose: ((event: { code?: number; reason?: string; wasClean: boolean }) => void) | null = null
  public onerror: ((event: unknown) => void) | null = null
  public onmessage: ((event: { data: unknown }) => void) | null = null
  public url: string

  constructor(url: string) {
    this.url = url
    MockWebSocket.instances.push(this)
  }

  close(code?: number, reason?: string) {
    this.readyState = MockWebSocket.CLOSED
    this.onclose?.({ code, reason, wasClean: true })
  }

  send(): void {}
}

let OriginalWebSocket: typeof WebSocket

beforeEach(() => {
  vi.useFakeTimers()
  MockWebSocket.instances = []
  OriginalWebSocket = globalThis.WebSocket as typeof WebSocket
  globalThis.WebSocket = MockWebSocket as unknown as typeof WebSocket
})

afterEach(() => {
  vi.useRealTimers()
  globalThis.WebSocket = OriginalWebSocket
})

describe('WebSocketService', () => {
  it('does not reconnect after reaching maxReconnectAttempts', () => {
    const service = new WebSocketService({ url: 'ws://test', maxReconnectAttempts: 2, reconnectInterval: 10 })

    service.connect()
    expect(MockWebSocket.instances.length).toBe(1)

    // First disconnection triggers first reconnection
    MockWebSocket.instances[0].close()
    vi.runAllTimers()
    expect(MockWebSocket.instances.length).toBe(2)

    // Second disconnection triggers second reconnection
    MockWebSocket.instances[1].close()
    vi.runAllTimers()
    expect(MockWebSocket.instances.length).toBe(3)

    // Third disconnection should not reconnect
    MockWebSocket.instances[2].close()
    vi.runAllTimers()
    expect(MockWebSocket.instances.length).toBe(3)
  })

  it('stops reconnecting after disconnect is called', () => {
    const service = new WebSocketService({ url: 'ws://test', maxReconnectAttempts: 5, reconnectInterval: 10 })

    service.connect()
    expect(MockWebSocket.instances.length).toBe(1)

    service.disconnect()
    vi.runAllTimers()
    expect(MockWebSocket.instances.length).toBe(1)
  })
})
