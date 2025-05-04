import type {IAlgorithmBiologyConfigBase} from "$lib/stores/algorithms.svelte";

export interface INSGAIIConfig extends IAlgorithmBiologyConfigBase {
  crossoverRate: number
  mutationRate: number
  mutationStrength: number
  tournamentSize: number
}

export const nsgaiiConfig = $state<INSGAIIConfig>({
  chromosome: 100,
  generation: 300,
  mutationStrength: 0.1,
  type: 'Biology',
  crossoverRate: 0.9,
  mutationRate: 0.1,
  tournamentSize: 10,
})
