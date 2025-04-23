
export interface ISizeConfig {
  AlphaSizePenalty: number
  PowerDifferencePenalty: number
  SmallLocations: string[],
  LargeFacilities: string[]
}


export const sizeConfig = $state<ISizeConfig>({
  AlphaSizePenalty: 20000,
  PowerDifferencePenalty: 1,
  LargeFacilities: [],
  SmallLocations: [],
})