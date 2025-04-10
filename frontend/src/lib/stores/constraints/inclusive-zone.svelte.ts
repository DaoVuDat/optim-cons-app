
export interface IZone {
  Name: string;
  BuildingNames: string;
  Size: number;
  Id: string;
}

export interface IInclusiveZoneConfig {
  AlphaInclusiveZonePenalty: number,
  PowerDifferencePenalty: number
  Zones: IZone[]
}


export const inclusiveZoneConfig = $state<IInclusiveZoneConfig>({
  AlphaInclusiveZonePenalty: 20000,
  PowerDifferencePenalty: 1,
  Zones: []
})