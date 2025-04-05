import type {IAlgorithmBiologyConfigBase} from "$lib/stores/algorithms.svelte";

export interface IGAConfig extends IAlgorithmBiologyConfigBase{
  mutationRate: number
  crossoverRate: number
}


export const gaConfig = $state<IGAConfig>({
  chromosome: {
    label: 'Chromosome',
    value: 100,
  },
  generation: {
    label: 'Generation',
    value: 300,
  },
  type: 'Biology',
  crossoverRate: 0.7,
  mutationRate: 0.1
})