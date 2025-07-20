import type { IntRange } from 'type-fest'

export namespace ArtNet {
  export type UniverseData = Record<Universe, DmxValue[]>
  export type Universe = IntRange<0, 0xffff>
  export type DmxChannel = IntRange<0, 511>
  export type DmxValue = IntRange<0, 255>

  export interface ArtDMXPacket {
    Sequence: number
    Physical: number
    SubUni: number
    Net: number
    Length: number
    Data: DmxValue[]
  }

  export interface ArtNetNode {
    IPAddress: string
    ShortName: string
    LongName: string
    NodeReport: string
    MacAddress: string
    LastSeen: string
  }
}
