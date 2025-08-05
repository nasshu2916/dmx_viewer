import React, { useEffect, useState, useCallback } from 'react'
import TimeDisplay from './TimeDisplay'

function formatTime(date: Date): string {
  return date.toLocaleTimeString('ja-JP', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
    timeZoneName: 'short',
  })
}

const TimeDisplayContainer: React.FC = () => {
  const [baseTime, setBaseTime] = useState<Date | null>(null)
  const [baseLocal, setBaseLocal] = useState<number | null>(null)
  const [timeText, setTime] = useState<string>(formatTime(new Date()))

  const getAdjustedTime = useCallback(() => {
    if (baseTime && baseLocal) {
      return new Date(baseTime.getTime() + (Date.now() - baseLocal))
    } else {
      return new Date()
    }
  }, [baseTime, baseLocal])

  const fetchServerTime = useCallback(async () => {
    const res = await fetch('/api/time')
    if (!res.ok) {
      console.warn('サーバから時刻を取得できませんでした')
      return
    }
    const data = await res.json()
    const ntpDate = new Date(data.datetime)
    setBaseTime(ntpDate)
    setBaseLocal(Date.now())
    setTime(formatTime(ntpDate))
  }, [])

  useEffect(() => {
    const timer = setInterval(() => {
      setTime(formatTime(getAdjustedTime()))
    }, 100)
    return () => clearInterval(timer)
  }, [getAdjustedTime])

  return <TimeDisplay formatTime={timeText} onReload={fetchServerTime} />
}

export default React.memo(TimeDisplayContainer)
