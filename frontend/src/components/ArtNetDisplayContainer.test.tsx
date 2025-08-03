import { describe, it, expect } from 'vitest'
import { calcColumns, getNextChannelByKey } from './ArtNetDisplayContainer'

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

describe('getNextChannelByKey', () => {
  const maxChannel = 511
  it('ArrowUp: moves up a row', () => {
    expect(getNextChannelByKey('ArrowUp', 20, 8, maxChannel)).toBe(12)
  })
  it('ArrowDown: moves down a row', () => {
    expect(getNextChannelByKey('ArrowDown', 20, 8, maxChannel)).toBe(28)
  })
  it('ArrowLeft: moves left a column', () => {
    expect(getNextChannelByKey('ArrowLeft', 20, 8, maxChannel)).toBe(19)
  })
  it('ArrowRight: moves right a column', () => {
    expect(getNextChannelByKey('ArrowRight', 20, 8, maxChannel)).toBe(21)
  })
  it('ArrowUp at top row returns null', () => {
    expect(getNextChannelByKey('ArrowUp', 3, 8, maxChannel)).toBeNull()
  })
  it('ArrowLeft at leftmost column returns null', () => {
    expect(getNextChannelByKey('ArrowLeft', 16, 8, maxChannel)).toBeNull()
  })
  it('ArrowDown at last row returns null', () => {
    expect(getNextChannelByKey('ArrowDown', 504, 8, maxChannel)).toBeNull()
  })
  it('ArrowRight at rightmost column returns null', () => {
    expect(getNextChannelByKey('ArrowRight', 23, 8, maxChannel)).toBeNull()
  })
  it('Unknown key returns null', () => {
    expect(getNextChannelByKey('a', 20, 8, maxChannel)).toBeNull()
  })
})
