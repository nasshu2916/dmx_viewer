import { useEffect, useRef, useState } from 'react'

export interface DmxHistoryPoint {
  value: number
  timestamp: number
}

/**
 * 選択中Universe/ChannelのDMX値ヒストリーを管理するフック
 * @param dmxValue 現在のDMX値
 * @param selectedKey 履歴をリセットするためのキー
 * @param maxLength 履歴の最大長
 */
export function useDmxHistory(dmxValue: number | null, selectedKey: string, maxLength = 100): DmxHistoryPoint[] {
  const [history, setHistory] = useState<DmxHistoryPoint[]>([])
  const prevKey = useRef<string | undefined>(selectedKey)

  useEffect(() => {
    if (!selectedKey) {
      setHistory([])
      prevKey.current = selectedKey
      return
    }

    if (prevKey.current !== selectedKey) {
      setHistory([])
      prevKey.current = selectedKey
    }

    if (dmxValue != null) {
      setHistory(prev => {
        const sliced = prev.length >= maxLength ? prev.slice(prev.length - maxLength + 1) : prev
        return [...sliced, { value: dmxValue, timestamp: Date.now() }]
      })
    }
  }, [dmxValue, selectedKey, maxLength])

  return history
}
