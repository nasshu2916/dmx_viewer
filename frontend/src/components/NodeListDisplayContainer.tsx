import React from 'react'
import { useSelection } from '@/contexts/SelectionContext'
import NodeListDisplay from './NodeListDisplay'
import { useWebSocket } from '@/contexts/WebSocketContext'
import type { ArtNet } from '@/types/artnet'

export type NodeListDisplayNode = {
  address: string
  info: ArtNet.ArtNetNode
  universes: ArtNet.Universe[]
  isUnknown: boolean
}

const NodeListDisplayContainer: React.FC = () => {
  const { setSelectedUniverse } = useSelection()
  const { artNetNodes, dmxData } = useWebSocket()
  const nodes = Array.isArray(artNetNodes) ? artNetNodes : []
  const receiveUniverseByNode = new Map<string, ArtNet.Universe[]>()
  for (const [address, universes] of Object.entries(dmxData)) {
    const universeNumbers: ArtNet.Universe[] = Object.keys(universes).map(Number) as ArtNet.Universe[]
    receiveUniverseByNode.set(address, universeNumbers)
  }
  // ノードが存在しないがdmxDataにだけ存在するアドレス
  const missingAddresses = Array.from(receiveUniverseByNode.keys()).filter(
    address => !nodes.some(node => node.IPAddress === address)
  )

  const displayNodes: NodeListDisplayNode[] = [
    ...nodes.map(node => ({
      address: node.IPAddress,
      info: node,
      universes: receiveUniverseByNode.get(node.IPAddress) || [],
      isUnknown: false,
    })),
    ...missingAddresses.map(address => ({
      address,
      info: invalidNode(address),
      universes: receiveUniverseByNode.get(address) || [],
      isUnknown: true,
    })),
  ]

  const handleSelectUniverses = (address: string, selected: ArtNet.Universe) => {
    setSelectedUniverse({ address, universe: selected })
  }
  return <NodeListDisplay nodes={displayNodes} onSelectUniverses={handleSelectUniverses} />
}

function invalidNode(address: string): ArtNet.ArtNetNode {
  return {
    IPAddress: address,
    ShortName: 'Unknown',
    LongName: 'Unknown Node',
    NodeReport: '',
    MacAddress: '00:00:00:00:00:00',
  }
}

export default NodeListDisplayContainer
