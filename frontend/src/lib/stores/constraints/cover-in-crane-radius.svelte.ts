

export interface ICoverInCraneRadiusConfig {
  AlphaCoverInCraneRadiusPenalty: number;
  PowerDifferencePenalty: number
}


export const coverInCraneRadiusConfig = $state<ICoverInCraneRadiusConfig>({
  AlphaCoverInCraneRadiusPenalty: 20000,
  PowerDifferencePenalty: 1
})