export * from './artnet'
export * from './websocket'

import type { ArtNet } from './artnet'

export type SelectedUniverse = {
  address: string
  universe: ArtNet.Universe
}
