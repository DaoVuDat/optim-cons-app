import type {IAlgorithmBiologyConfigBase} from "$lib/stores/algorithms.svelte";

export interface IGAConfig extends IAlgorithmBiologyConfigBase {
  mutationRate: number
  crossoverRate: number
  elitismCount: number
}


export const gaConfig = $state<IGAConfig>({
  chromosome: 100,
  generation: 300,
  type: 'Biology',
  crossoverRate: 0.7,
  mutationRate: 0.1,
  elitismCount: 5,
})