import { describe, it, expect } from 'vitest'
import { getUniverse } from './artnet'
import type { ArtNet } from '@/types/artnet'

describe('getUniverse', () => {
  it('calculates universe from Net and SubUni', () => {
    const packet: ArtNet.ArtDMXPacket = {
      Net: 1,
      SubUni: 2,
      Data: new Array(512).fill(0) as ArtNet.DmxValue[],
      Length: 512,
      Sequence: 0,
      Physical: 0,
      SourceIP: '127.0.0.1',
    }
    expect(getUniverse(packet)).toBe(258) // (1 << 8) | 2 = 256 + 2
  })

  it('returns 0 for Net=0, SubUni=0', () => {
    const packet: ArtNet.ArtDMXPacket = {
      Net: 0,
      SubUni: 0,
      Data: new Array(512).fill(0) as ArtNet.DmxValue[],
      Length: 512,
      Sequence: 0,
      Physical: 0,
      SourceIP: '127.0.0.1',
    }
    expect(getUniverse(packet)).toBe(0)
  })

  it('returns max value for Net=255, SubUni=255', () => {
    const packet: ArtNet.ArtDMXPacket = {
      Net: 255,
      SubUni: 255,
      Data: new Array(512).fill(0) as ArtNet.DmxValue[],
      Length: 512,
      Sequence: 0,
      Physical: 0,
      SourceIP: '127.0.0.1',
    }
    expect(getUniverse(packet)).toBe(65535) // (255 << 8) | 255 = 65280 + 255
  })
})
