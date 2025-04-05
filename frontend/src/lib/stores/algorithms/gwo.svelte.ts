import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IGWOConfig extends IAlgorithmSwarmConfigBase {
}


export const   gwoConfig = $state<IGWOConfig>({
  iterations: {
    label: 'Iterations',
    value: 300,
  },
  population: {
    label: 'Population',
    value: 100,
  },
  type: 'Swarm',
})

