import React from 'react'
import ArtNetDisplay from './ArtNetDisplay'
import { useWebSocket } from '@/contexts/WebSocketContext'

import type { ArtNet } from '@/types/artnet'

/**
 * columns（2のN乗）を計算するユーティリティ関数
 * @param containerWidth - コンテナの幅(px)
 * @param cellWidth - セルの最小幅(px)
 * @param minColumns - 最小カラム数（デフォルト1）
 * @param maxColumns - 最大カラム数（デフォルト32）
 * @returns 2のN乗のカラム数（minColumns <= columns <= maxColumns）
 */
export function calcColumns(containerWidth: number, cellWidth = 46, minColumns = 1, maxColumns = 32): number {
  let cols = Math.max(minColumns, Math.min(maxColumns, Math.floor(containerWidth / cellWidth)))
  // 2のN乗に切り捨て
  cols = Math.pow(2, Math.floor(Math.log2(cols)))
  if (cols < minColumns) cols = minColumns
  if (cols > maxColumns) cols = maxColumns
  return cols
}

/**
 * キー移動ロジック: 新しいチャンネル番号を返す
 * @param key - 押されたキー（ArrowUp, ArrowDown, ArrowLeft, ArrowRight）
 * @param selectedChannel - 現在のチャンネル番号
 * @param columns - 1行あたりのカラム数
 * @param maxChannel - 最大チャンネル番号
 * @returns 新しいチャンネル番号（変化なしの場合は同じ値）
 */
/**
 * キー移動ロジック: 新しいチャンネル番号を返す
 * @param key - 押されたキー（ArrowUp, ArrowDown, ArrowLeft, ArrowRight）
 * @param selectedChannel - 現在のチャンネル番号
 * @param columns - 1行あたりのカラム数
 * @param maxChannel - 最大チャンネル番号
 * @returns 新しいチャンネル番号（変化がなければnull）
 * @internal テスト以外でimportしないこと
 */
export function getNextChannelByKey(
  key: string,
  selectedChannel: number,
  columns: number,
  maxChannel: number
): ArtNet.DmxChannel | null {
  const row = Math.floor(selectedChannel / columns)
  const col = selectedChannel % columns
  let newRow = row
  let newCol = col
  const rows = Math.ceil((maxChannel + 1) / columns)
  switch (key) {
    case 'ArrowUp':
      newRow = Math.max(0, row - 1)
      break
    case 'ArrowDown':
      newRow = Math.min(rows - 1, row + 1)
      break
    case 'ArrowLeft':
      newCol = Math.max(0, col - 1)
      break
    case 'ArrowRight':
      newCol = Math.min(columns - 1, col + 1)
      break
    default:
      return null
  }
  let newChannel = newRow * columns + newCol
  if (newChannel > maxChannel) newChannel = maxChannel
  if (newChannel === selectedChannel) return null
  return newChannel as ArtNet.DmxChannel
}

interface ArtNetDisplayContainerProps {
  displayUniverse?: [string, number]
  selectedChannel: ArtNet.DmxChannel | null
  onSelectChannel: (channel: ArtNet.DmxChannel) => void
}

const ArtNetDisplayContainer: React.FC<ArtNetDisplayContainerProps> = ({
  displayUniverse,
  selectedChannel,
  onSelectChannel,
}) => {
  // 横並び数（columns）をContainerで管理
  const [columns, setColumns] = React.useState(16)
  // eslint-disable-next-line no-undef
  const containerRef = React.useRef<HTMLDivElement>(null)

  const calcColumnsCallback = React.useCallback(() => {
    const containerWidth = containerRef.current?.clientWidth ?? window.innerWidth
    const cols = calcColumns(containerWidth)
    setColumns(cols)
  }, [])

  React.useEffect(() => {
    calcColumnsCallback()
    window.addEventListener('resize', calcColumnsCallback)
    return () => window.removeEventListener('resize', calcColumnsCallback)
  }, [calcColumnsCallback])
  const { dmxData } = useWebSocket()
  const dmxDataForDisplay = Object.fromEntries(
    Object.entries(dmxData).map(([address, universes]) => [
      address,
      Object.fromEntries(
        Object.entries(universes).map(([universe, obj]) => [
          Number(universe),
          { data: obj.data, receivedAt: obj.receivedAt },
        ])
      ),
    ])
  )

  // キー移動
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (selectedChannel === null) return
    const maxChannel = 511
    const newChannel = getNextChannelByKey(e.key, selectedChannel, columns, maxChannel)
    if (newChannel !== null) {
      e.preventDefault()
      onSelectChannel(newChannel)
    }
  }

  return (
    <div
      className="outline-none focus:ring-1 focus:ring-dmx-accent"
      ref={containerRef}
      style={{ outline: 'none' }}
      tabIndex={0}
      onKeyDown={handleKeyDown}
    >
      <ArtNetDisplay
        columns={columns}
        displayUniverse={displayUniverse}
        dmxData={dmxDataForDisplay}
        selectedChannel={selectedChannel}
        onSelectChannel={onSelectChannel}
      />
    </div>
  )
}

export default ArtNetDisplayContainer
