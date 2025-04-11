
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
  Zones: [
    {
      Size: 20,
      Name: 'TF13',
      Id: Math.random().toString(),
      BuildingNames: 'TF7'
    },
    {
      Size: 15,
      Name: 'TF13',
      Id: Math.random().toString(),
      BuildingNames: 'TF1 TF2'
    }
  ]
})