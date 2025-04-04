import type {IAlgorithmBiologyConfigBase} from "$lib/stores/algorithms.svelte";

export interface IGAConfig extends IAlgorithmBiologyConfigBase{
  mutationRate: number
  crossoverRate: number
}
