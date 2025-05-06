import type {IAlgorithmBiologyConfigBase} from "$lib/stores/algorithms.svelte";

export interface INSGAIIConfig extends IAlgorithmBiologyConfigBase {
  crossoverRate: number
  mutationRate: number
  mutationStrength: number
  sigma: number
}

export const nsgaiiConfig = $state<INSGAIIConfig>({
  chromosome: 100,
  generation: 300,
  crossoverRate: 0.7,
  mutationRate: 0.4,
  mutationStrength: 0.01,
  sigma: 0.1,
  type: 'Biology',
})
