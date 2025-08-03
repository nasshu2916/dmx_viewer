import { describe, it, expect } from 'vitest'
import { calcColumns } from './artnetDisplayUtils'

describe('calcColumns', () => {
  it('returns 2 when containerWidth is just enough for 2 cells', () => {
    expect(calcColumns(96, 48)).toBe(2)
  })
  it('returns 4 when containerWidth is just enough for 4 cells', () => {
    expect(calcColumns(192, 48)).toBe(4)
  })
  it('returns 8 when containerWidth is just enough for 8 cells', () => {
    expect(calcColumns(384, 48)).toBe(8)
  })
  it('returns 32 when containerWidth is very large', () => {
    expect(calcColumns(2000, 48)).toBe(32)
  })
  it('returns 1 when containerWidth is very small', () => {
    expect(calcColumns(10, 48)).toBe(1)
  })
  it('respects minColumns and maxColumns', () => {
    expect(calcColumns(10, 48, 4, 8)).toBe(4)
    expect(calcColumns(2000, 48, 4, 8)).toBe(8)
  })
  it('returns 2^N only', () => {
    expect([1, 2, 4, 8, 16, 32]).toContain(calcColumns(1000, 48))
  })
})
