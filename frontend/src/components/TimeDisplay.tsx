import React, { useEffect, useState } from 'react'

const TimeDisplay: React.FC = () => {
  const [baseTime, setBaseTime] = useState<Date | null>(null)
  const [baseLocal, setBaseLocal] = useState<number | null>(null)
  const [time, setTime] = useState(new Date())

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('ja-JP', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
      timeZoneName: 'short',
    })
  }

  const getAdjustedTime = () => {
    if (baseTime && baseLocal) {
      return new Date(baseTime.getTime() + (Date.now() - baseLocal))
    } else {
      return new Date()
    }
  }

  const fetchServerTime = async () => {
    const res = await fetch('/api/ntp_time')
    if (!res.ok) {
      console.warn('NTPサーバから時刻を取得できませんでした')
      return
    }
    const data = await res.json()
    const ntpDate = new Date(data.datetime)
    setBaseTime(ntpDate)
    setBaseLocal(Date.now())
    setTime(ntpDate)
  }

  useEffect(() => {
    const timer = setInterval(() => {
      setTime(getAdjustedTime())
    }, 100)
    return () => clearInterval(timer)
  }, [baseTime, baseLocal])

  return (
    <div className="flex items-center space-x-2 p-2 text-dmx-text-light">
      <span className="m-0 text-lg">{formatTime(time)}</span>
      <button
        className="flex cursor-pointer rounded bg-transparent p-2 text-dmx-text-light hover:bg-white/10"
        onClick={fetchServerTime}
      >
        <svg
          className="lucide lucide-rotate-cw h-4 w-4"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8" />
          <path d="M21 3v5h-5" />
        </svg>
      </button>
    </div>
  )
}

export default TimeDisplay
