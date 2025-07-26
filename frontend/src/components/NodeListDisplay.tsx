import React from 'react'
import type { ArtNet } from '@/types/artnet'

interface NodeListDisplayProps {
  artNetNodes: ArtNet.ArtNetNode[]
}

const NodeListDisplay: React.FC<NodeListDisplayProps> = ({ artNetNodes }) => {
  // Ensure artNetNodes is always an array
  const nodes = Array.isArray(artNetNodes) ? artNetNodes : []

  return (
    <div className="p-4">
      <h2 className="mb-4 text-xl font-bold">ArtNet Nodes</h2>
      {nodes.length === 0 ? (
        <p>No ArtNet nodes discovered yet.</p>
      ) : (
        <ul>
          {nodes.map((node: ArtNet.ArtNetNode) => (
            <li className="mb-2 rounded border border-gray-700 p-2" key={node.IPAddress}>
              <p className="font-semibold">{node.ShortName || 'Unknown Node'}</p>
              <p className="text-sm text-gray-400">IP: {node.IPAddress}</p>
              <p className="text-sm text-gray-400">MAC: {node.MacAddress}</p>
              <p className="text-sm text-gray-400">Last Seen: {new Date(node.LastSeen).toLocaleString()}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default NodeListDisplay
