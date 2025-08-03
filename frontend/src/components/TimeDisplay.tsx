import React from 'react'

interface TimeDisplayProps {
  time: Date
  onReload: () => void
}

const TimeDisplay: React.FC<TimeDisplayProps> = ({ time, onReload }) => {
  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('ja-JP', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
      timeZoneName: 'short',
    })
  }

  return (
    <div className="flex items-center space-x-2 p-2 text-dmx-text-light">
      <span className="m-0 text-lg">{formatTime(time)}</span>
      <button
        className="flex cursor-pointer rounded bg-transparent p-2 text-dmx-text-light hover:bg-white/10"
        onClick={onReload}
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
