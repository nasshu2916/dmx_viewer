import type { IntRange } from 'type-fest'

export namespace ArtNet {
  export type UniverseData = Record<Universe, DmxValue[]>
  /** Universe: 0〜65535 の範囲の number */
  export type Universe = number
  export type DmxChannel = IntRange<0, 511>
  export type DmxValue = IntRange<0, 255>

  export interface ArtDMXPacket {
    Sequence: number
    Physical: number
    SubUni: number
    Net: number
    Length: number
    Data: DmxValue[]
    SourceIP: string
  }

  export interface ArtNetNode {
    IPAddress: string
    ShortName: string
    LongName: string
    NodeReport: string
    MacAddress: string
    LastSeen?: string
  }
}
