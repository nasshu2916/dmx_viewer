import React from 'react'

interface TimeDisplayProps {
  formatTime: string
  onReload: () => void
}

const TimeDisplay: React.FC<TimeDisplayProps> = React.memo(({ formatTime, onReload }) => {
  return (
    <div className="flex items-center space-x-1 text-dmx-text-light">
      <span className="m-0 text-xs md:text-sm leading-none">{formatTime}</span>
      <button
        className="flex cursor-pointer rounded bg-transparent p-1 text-dmx-text-light hover:bg-white/10"
        onClick={onReload}
      >
        <svg
          className="lucide lucide-rotate-cw h-3 w-3 md:h-4 md:w-4"
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
})

export default TimeDisplay
