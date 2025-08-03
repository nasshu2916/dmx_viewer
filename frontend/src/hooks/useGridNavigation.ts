import type React from 'react'
import { useCallback } from 'react'

/**
 * 汎用的なグリッドナビゲーションフック
 * @param currentIndex - 現在のインデックス
 * @param rowCount - 行数
 * @param colCount - 列数
 * @param onMove - 移動時コールバック
 * @param isCellValid - (index: number) => boolean 有効セル判定（省略可）
 */
export function useGridNavigation({
  currentIndex,
  rowCount,
  colCount,
  onMove,
  isCellValid,
}: {
  currentIndex: number
  rowCount: number
  colCount: number
  onMove: (nextIndex: number) => void
  isCellValid?: (index: number) => boolean
}) {
  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      let nextIndex = currentIndex
      const row = Math.floor(currentIndex / colCount)
      const col = currentIndex % colCount
      let newRow = row
      let newCol = col
      switch (e.key) {
        case 'ArrowUp':
          newRow = Math.max(0, row - 1)
          break
        case 'ArrowDown':
          newRow = Math.min(rowCount - 1, row + 1)
          break
        case 'ArrowLeft':
          newCol = Math.max(0, col - 1)
          break
        case 'ArrowRight':
          newCol = Math.min(colCount - 1, col + 1)
          break
        default:
          return
      }
      nextIndex = newRow * colCount + newCol
      if (isCellValid && !isCellValid(nextIndex)) return
      if (nextIndex !== currentIndex) {
        e.preventDefault()
        onMove(nextIndex)
      }
    },
    [currentIndex, rowCount, colCount, onMove, isCellValid]
  )
  return { handleKeyDown }
}
