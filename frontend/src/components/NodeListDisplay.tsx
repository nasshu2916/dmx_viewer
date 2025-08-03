import React, { useState } from 'react'
import type { ArtNet } from '@/types/artnet'

interface NodeListDisplayProps {
  artNetNodes: ArtNet.ArtNetNode[]
  dmxData: Record<string, Record<number, ArtNet.DmxValue[]>>
  onSelectUniverses: (address: string, selected: number) => void
}

interface NodeUniverseListProps {
  address: string
  universes: number[]
  onSelectUniverses: (address: string, selected: number) => void
}

const NodeUniverseList: React.FC<NodeUniverseListProps> = ({ address, universes, onSelectUniverses }) => {
  const [selectedUniverse, setSelectedUniverse] = useState<[string, number] | undefined>(undefined)
  const handleRadioChange = (address: string, universe: number) => {
    setSelectedUniverse([address, universe])
    onSelectUniverses(address, universe)
  }

  return (
    <div className="mb-2">
      <p className="text-sm text-gray-500">Universes for {address}:</p>
      {universes.length > 0 ? (
        <div className="flex flex-col gap-2">
          {universes.map(universe => (
            <label className="flex items-center text-dmx-text-light" key={universe}>
              <input
                checked={selectedUniverse?.[0] === address && selectedUniverse?.[1] === universe}
                className="form-radio h-4 w-4 text-dmx-accent focus:ring-dmx-accent"
                type="radio"
                onChange={() => handleRadioChange(address, universe)}
              />
              <span className="ml-2">Universe {universe}</span>
            </label>
          ))}
        </div>
      ) : (
        <p className="text-sm text-gray-500">No universes received.</p>
      )}
    </div>
  )
}

const NodeInfo: React.FC<{ node: ArtNet.ArtNetNode }> = ({ node }) => {
  const lastSeen = node.LastSeen ? new Date(node.LastSeen).toLocaleString() : 'Unknown'

  return (
    <div>
      <p className="font-semibold">{node.ShortName || 'Unknown Node'}</p>
      <p className="text-sm text-gray-400">IP: {node.IPAddress}</p>
      <p className="text-sm text-gray-400">MAC: {node.MacAddress}</p>
      <p className="text-sm text-gray-400">Last Seen: {lastSeen}</p>
    </div>
  )
}

const NodeListDisplay: React.FC<NodeListDisplayProps> = ({ artNetNodes, dmxData, onSelectUniverses }) => {
  // Ensure artNetNodes is always an array
  const nodes = Array.isArray(artNetNodes) ? artNetNodes : []
  const receiveUniverseByNode = new Map<string, number[]>()
  for (const [address, universes] of Object.entries(dmxData)) {
    const universeNumbers = Object.keys(universes).map(Number)
    receiveUniverseByNode.set(address, universeNumbers)
  }
  // node が存在しない受信した Universe
  const missingUniversesByNode = Array.from(receiveUniverseByNode.keys()).filter(
    address => !nodes.some(node => node.IPAddress === address)
  )

  return (
    <div className="p-4">
      <h2 className="mb-4 text-xl font-bold">ArtNet Nodes</h2>
      <ul>
        {nodes.map((node: ArtNet.ArtNetNode) => (
          <li className="mb-2 rounded border border-gray-700 p-2" key={node.IPAddress}>
            <NodeInfo key={node.IPAddress} node={node} />
            <NodeUniverseList
              address={node.IPAddress}
              universes={receiveUniverseByNode.get(node.IPAddress) || []}
              onSelectUniverses={onSelectUniverses}
            />
          </li>
        ))}
        {missingUniversesByNode.map(address => (
          <li className="mb-2 rounded border border-gray-700 p-2" key={address}>
            <NodeInfo key={address} node={invalidAddressNode(address)} />
            <NodeUniverseList
              address={address}
              universes={receiveUniverseByNode.get(address) || []}
              onSelectUniverses={onSelectUniverses}
            />
          </li>
        ))}
      </ul>
    </div>
  )
}

function invalidAddressNode(address: string): ArtNet.ArtNetNode {
  return {
    IPAddress: address,
    ShortName: 'Unknown Node',
    LongName: 'Unknown Node',
    NodeReport: '',
    MacAddress: '00:00:00:00:00:00',
  }
}

export default NodeListDisplay
