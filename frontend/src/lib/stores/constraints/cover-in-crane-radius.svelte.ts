import type {ISelectedCraneWithId} from "$lib/stores/objectives";


export interface ICoverInCraneRadiusConfig {
  CraneLocations: ISelectedCraneWithId[];
  AlphaCoverInCraneRadiusPenalty: number;
  PowerDifferencePenalty: number
}


export const coverInCraneRadiusConfig = $state<ICoverInCraneRadiusConfig>({
  AlphaCoverInCraneRadiusPenalty: 20000,
  PowerDifferencePenalty: 1,
  CraneLocations: [],
})