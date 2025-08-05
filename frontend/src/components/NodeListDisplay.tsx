import React from 'react'
import type { ArtNet } from '@/types/artnet'
import type { NodeListDisplayNode } from './NodeListDisplayContainer'
import { useSelection } from '@/contexts/SelectionContext'

interface NodeListDisplayProps {
  nodes: NodeListDisplayNode[]
}

interface NodeUniverseListProps {
  address: string
  universes: ArtNet.Universe[]
}

const NodeUniverseList: React.FC<NodeUniverseListProps> = ({ address, universes }) => {
  const { selectedUniverse, setSelectedUniverse } = useSelection()

  const handleRadioChange = (address: string, universe: ArtNet.Universe) => {
    setSelectedUniverse({ address: address, universe: universe })
  }

  return (
    <div className="my-2">
      {universes.length > 0 ? (
        <div className="flex flex-row gap-2">
          {universes.map(universe => {
            const isSelected = selectedUniverse?.address === address && selectedUniverse?.universe === universe
            return (
              <button
                className={`ml-0 rounded border-2 px-4 py-2 transition-colors focus:outline-none ${isSelected ? 'border-dmx-accent bg-dmx-accent/10 text-dmx-accent' : 'border border-gray-600 bg-transparent text-dmx-text-light hover:bg-dmx-accent/5'}`}
                key={universe}
                type="button"
                onClick={() => handleRadioChange(address, universe)}
              >
                {universe}
              </button>
            )
          })}
        </div>
      ) : (
        <p className="text-sm text-gray-500">No universes received.</p>
      )}
    </div>
  )
}

const NodeInfo: React.FC<{ node: ArtNet.ArtNetNode }> = React.memo(({ node }) => {
  const lastSeen = node.LastSeen ? new Date(node.LastSeen).toLocaleString() : 'Unknown'

  return (
    <div>
      <p className="font-semibold">{node.ShortName || 'Unknown Node'}</p>
      <p className="text-sm text-gray-400">IP: {node.IPAddress}</p>
      <p className="text-sm text-gray-400">MAC: {node.MacAddress}</p>
      <p className="text-sm text-gray-400">Last Seen: {lastSeen}</p>
    </div>
  )
})

NodeInfo.displayName = 'NodeInfo'

const NodeListDisplay: React.FC<NodeListDisplayProps> = ({ nodes }) => {
  return (
    <div className="p-4">
      <h2 className="mb-4 text-xl font-bold">ArtNet Nodes</h2>
      <ul>
        {nodes.map(node => (
          <li className="mb-2 rounded border border-gray-700 p-2" key={node.address}>
            <NodeInfo node={node.info} />
            <NodeUniverseList address={node.address} universes={node.universes} />
          </li>
        ))}
      </ul>
    </div>
  )
}

export default NodeListDisplay
