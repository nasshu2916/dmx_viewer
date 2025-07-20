import React, { useState } from 'react'

interface UniverseSelectorProps {
  availableUniverses: number[]
  onSelectUniverses: (selected: number) => void
}

const UniverseSelector: React.FC<UniverseSelectorProps> = ({ availableUniverses, onSelectUniverses }) => {
  const [selectedUniverse, setSelectedUniverse] = useState<number | undefined>(undefined)

  const handleRadioChange = (universe: number) => {
    setSelectedUniverse(universe)
    onSelectUniverses(universe)
  }

  return (
    <div className="rounded-lg bg-dmx-light-bg p-4 shadow-lg">
      <h3 className="mb-4 text-lg font-bold text-dmx-text-light">Select Universe</h3>
      <div className="flex flex-col gap-2">
        {availableUniverses.map(universe => (
          <label className="flex items-center text-dmx-text-light" key={universe}>
            <input
              checked={selectedUniverse === universe}
              className="form-radio h-4 w-4 text-dmx-accent focus:ring-dmx-accent"
              type="radio"
              onChange={() => handleRadioChange(universe)}
            />
            <span className="ml-2">Universe {universe}</span>
          </label>
        ))}
      </div>
    </div>
  )
}

export default UniverseSelector
