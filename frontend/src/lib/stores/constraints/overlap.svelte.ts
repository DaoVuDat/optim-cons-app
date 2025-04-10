

export interface IOverlapConfig {
  AlphaOverLapPenalty: number
  PowerDifferencePenalty: number
}


export const overlapConfig = $state<IOverlapConfig>({
  AlphaOverLapPenalty: 20000,
  PowerDifferencePenalty: 1,
})