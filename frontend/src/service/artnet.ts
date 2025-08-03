import type { ArtNet } from '@/types/artnet'

/**
 * ArtDMXPacket から Universe を取得する
 * Universe = (Net << 8) | SubUni
 */
export function getUniverse(packet: ArtNet.ArtDMXPacket): ArtNet.Universe {
  return (packet.Net << 8) | packet.SubUni
}
