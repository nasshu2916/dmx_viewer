import React, { useMemo, useCallback, useRef } from 'react'
import UniverseTable from './UniverseTable'
import { useWebSocket } from '@/contexts/WebSocketContext'
import { useGridNavigation } from '@/hooks/useGridNavigation'
import { useSelectionStore } from '@/stores/selectionStore'
import { calcColumns } from './artnetDisplayUtils'

import type { ArtNet } from '@/types/artnet'

/**
 * ArtNet DMXデータを表示するコンテナコンポーネント
 *
 * 機能:
 * - レスポンシブなカラム表示
 * - キーボードナビゲーション（矢印キー）
 * - WebSocketからのリアルタイムDMXデータ表示
 */
const ArtNetDisplayContainer: React.FC = () => {
  const { selectedUniverse: displayUniverse, selectedChannel, setSelectedChannel } = useSelectionStore()
  // レスポンシブなカラム数の管理（汎用hooks）
  // eslint-disable-next-line no-undef
  const containerRef = useRef<HTMLDivElement>(null)
  const [columns, setColumns] = React.useState(16)
  React.useEffect(() => {
    const updateColumns = () => {
      const containerWidth = containerRef.current?.clientWidth ?? window.innerWidth
      setColumns(calcColumns(containerWidth))
    }
    updateColumns()
    window.addEventListener('resize', updateColumns)
    return () => window.removeEventListener('resize', updateColumns)
  }, [])

  // WebSocketからDMXデータを取得
  const { dmxData } = useWebSocket()

  // DMXデータの変換とフィルタリング
  const { filteredData, universe, maxChannel } = useMemo(() => {
    // ユニバース情報を取得
    const addr = displayUniverse?.address ?? 'Unknown'
    const univ = displayUniverse?.universe ?? 0

    // DMXデータを変換
    const transformedData = Object.fromEntries(
      Object.entries(dmxData).map(([address, universes]) => [
        address,
        Object.fromEntries(
          Object.entries(universes as Record<string, { data: ArtNet.DmxValue[]; receivedAt: Date }>).map(
            ([universe, obj]) => [Number(universe), { data: obj.data, receivedAt: obj.receivedAt }]
          )
        ),
      ])
    )

    // フィルタされたデータを取得
    const filtered = transformedData[addr]?.[univ]

    return {
      filteredData: filtered,
      address: addr,
      universe: univ,
      maxChannel: filtered ? filtered.data.length - 1 : 0,
    }
  }, [dmxData, displayUniverse])

  // キーボードナビゲーションの処理（汎用hooks）
  const { handleKeyDown } = useGridNavigation({
    currentIndex: selectedChannel ?? 0,
    rowCount: Math.ceil((maxChannel + 1) / columns),
    colCount: columns,
    onMove: idx => setSelectedChannel(idx as ArtNet.DmxChannel),
    isCellValid: idx => idx >= 0 && idx <= maxChannel,
  })

  const renderWithDmxData = useCallback(
    (data: { data: ArtNet.DmxValue[]; receivedAt: Date }) => (
      <UniverseTable
        columns={columns}
        data={data.data}
        receivedAt={data.receivedAt}
        selectedChannel={selectedChannel}
        universe={universe}
        onSelectChannel={setSelectedChannel}
      />
    ),
    [columns, selectedChannel, universe, setSelectedChannel]
  )

  const renderWithoutDmxData = useCallback(() => <p className="text-dmx-text-light">Waiting for ArtNet data...</p>, [])

  return (
    <div
      className="outline-none focus:ring-1 focus:ring-dmx-accent"
      ref={containerRef}
      style={{ outline: 'none' }}
      tabIndex={0}
      onKeyDown={handleKeyDown}
    >
      {filteredData === undefined ? renderWithoutDmxData() : renderWithDmxData(filteredData)}
    </div>
  )
}

export default ArtNetDisplayContainer
