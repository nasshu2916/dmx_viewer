import React from 'react'
import NodeListDisplay from './NodeListDisplay'
import type { ArtNet } from '@/types/artnet'
import { useArtNetStore } from '@/stores'

export type NodeListDisplayNode = {
  address: string
  info: ArtNet.ArtNetNode
  universes: ArtNet.Universe[]
  isUnknown: boolean
}

const NodeListDisplayContainer: React.FC = () => {
  const { artNetNodes, dmxData } = useArtNetStore()
  const receiveUniverseByNode = new Map<string, ArtNet.Universe[]>()
  for (const [address, universes] of Object.entries(dmxData)) {
    const universeNumbers: ArtNet.Universe[] = Object.keys(universes).map(Number) as ArtNet.Universe[]
    receiveUniverseByNode.set(address, universeNumbers)
  }

  const displayNodes: NodeListDisplayNode[] = React.useMemo(() => {
    // ノードが存在しないがdmxDataにだけ存在するアドレス
    const missingAddresses = Array.from(receiveUniverseByNode.keys()).filter(
      address => !artNetNodes.some(node => node.IPAddress === address)
    )
    return [
      ...artNetNodes.map(node => ({
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
  }, [artNetNodes, receiveUniverseByNode])

  return <NodeListDisplay nodes={displayNodes} />
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

export default React.memo(NodeListDisplayContainer)
