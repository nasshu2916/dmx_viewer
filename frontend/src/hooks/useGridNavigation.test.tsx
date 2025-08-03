import type React from 'react'
import { renderHook, act } from '@testing-library/react'
import { useGridNavigation } from './useGridNavigation'
import { describe, it, expect, vi } from 'vitest'

describe('useGridNavigation', () => {
  const setup = (currentIndex: number, rowCount: number, colCount: number, isCellValid?: (idx: number) => boolean) => {
    const onMove = vi.fn()
    const { result } = renderHook(() => useGridNavigation({ currentIndex, rowCount, colCount, onMove, isCellValid }))
    return { handleKeyDown: result.current.handleKeyDown, onMove }
  }

  it('ArrowUp/Down/Left/Rightで正しく移動する', () => {
    const { handleKeyDown, onMove } = setup(5, 4, 4)
    // 5: row=1, col=1
    act(() => {
      handleKeyDown({ key: 'ArrowUp', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).toHaveBeenCalledWith(1) // 上
    act(() => {
      handleKeyDown({ key: 'ArrowDown', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).toHaveBeenCalledWith(9) // 下
    act(() => {
      handleKeyDown({ key: 'ArrowLeft', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).toHaveBeenCalledWith(4) // 左
    act(() => {
      handleKeyDown({ key: 'ArrowRight', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).toHaveBeenCalledWith(6) // 右
  })

  it('端で押しても移動しない', () => {
    const { handleKeyDown, onMove } = setup(0, 4, 4)
    act(() => {
      handleKeyDown({ key: 'ArrowUp', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).not.toHaveBeenCalled()
    act(() => {
      handleKeyDown({ key: 'ArrowLeft', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).not.toHaveBeenCalled()
  })

  it('isCellValidでfalseのセルには移動しない', () => {
    const { handleKeyDown, onMove } = setup(5, 4, 4, idx => idx !== 1)
    act(() => {
      handleKeyDown({ key: 'ArrowUp', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).not.toHaveBeenCalled() // 1は無効
  })

  it('無関係なキーでは何も起きない', () => {
    const { handleKeyDown, onMove } = setup(5, 4, 4)
    act(() => {
      handleKeyDown({ key: 'a', preventDefault: vi.fn() } as unknown as React.KeyboardEvent)
    })
    expect(onMove).not.toHaveBeenCalled()
  })
})
