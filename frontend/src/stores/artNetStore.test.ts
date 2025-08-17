import { act } from '@testing-library/react'
import { describe, test, expect, beforeEach } from 'vitest'
import { useArtNetStore } from './artNetStore'
import type { ServerMessage } from '@/types/websocket'

// Helper to get current state
const getState = () => useArtNetStore.getState()

const makeMessage = (i: number): ServerMessage => ({
  Type: 'info',
  Message: `msg-${i}`,
  Timestamp: i, // strictly increasing
})

describe('artNetStore serverMessages cap', () => {
  beforeEach(() => {
    act(() => {
      getState().clearData()
    })
  })

  test('keeps messages in order and caps at 200 newest', () => {
    act(() => {
      for (let i = 1; i <= 205; i++) {
        getState().addServerMessage(makeMessage(i))
      }
    })

    const { serverMessages } = getState()
    expect(serverMessages.length).toBe(200)
    // Should contain messages 6..205
    expect(serverMessages[0].Message).toBe('msg-6')
    expect(serverMessages[0].Timestamp).toBe(6)
    expect(serverMessages[199].Message).toBe('msg-205')
  })

  test('ignores duplicate timestamp at tail', () => {
    act(() => {
      getState().addServerMessage(makeMessage(1))
      getState().addServerMessage(makeMessage(1)) // duplicate timestamp
      getState().addServerMessage(makeMessage(2))
    })
    const { serverMessages } = getState()
    expect(serverMessages.length).toBe(2)
    expect(serverMessages[0].Timestamp).toBe(1)
    expect(serverMessages[1].Timestamp).toBe(2)
  })
})
