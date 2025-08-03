import React, { useEffect, useState, useCallback } from 'react'
import TimeDisplay from './TimeDisplay'

const TimeDisplayContainer: React.FC = () => {
  const [baseTime, setBaseTime] = useState<Date | null>(null)
  const [baseLocal, setBaseLocal] = useState<number | null>(null)
  const [time, setTime] = useState(new Date())

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
    setTime(ntpDate)
  }, [])

  useEffect(() => {
    const timer = setInterval(() => {
      setTime(getAdjustedTime())
    }, 100)
    return () => clearInterval(timer)
  }, [getAdjustedTime])

  return <TimeDisplay time={time} onReload={fetchServerTime} />
}

export default TimeDisplayContainer
